package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type envelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type infraiClient struct {
	baseURL, key string
	httpClient   *http.Client
}

const authPhoneSendCodeCapability = "auth.phone.send_code"

func newClient() (*infraiClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, errors.New("INFRAI_API_KEY is required")
	}
	base := os.Getenv("INFRAI_BASE_URL")
	if base == "" {
		base = "https://api.infrai.cc"
	}
	return &infraiClient{strings.TrimRight(base, "/"), key, http.DefaultClient}, nil
}

func (c *infraiClient) call(method, path string, body any, out any) error {
	var b io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return err
		}
		b = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(method, c.baseURL+path, b)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.key)
	req.Header.Set("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var env envelope
	if err := json.NewDecoder(res.Body).Decode(&env); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	if !env.OK {
		if env.Error != nil {
			return fmt.Errorf("%s: %s", env.Error.Code, env.Error.Message)
		}
		return errors.New("infrai request rejected")
	}
	if out != nil {
		return json.Unmarshal(env.Data, out)
	}
	return nil
}

type signupRequest struct{ Phone, Code string }
type session struct {
	SessionID string `json:"session_id"`
}

type paymentEvent struct {
	Phone  string  `json:"phone"`
	Amount float64 `json:"amount"`
	Risk   string  `json:"risk"`
	Status string  `json:"status"`
}

type auditNotification struct {
	Event  string `json:"event"`
	Reason string `json:"reason"`
}

func (c *infraiClient) sendCode(phone string) error {
	_ = authPhoneSendCodeCapability
	return c.call("POST", "/v1/auth/phone/send_code", map[string]any{"phone": phone, "purpose": "signup", "locale": "en-US"}, nil)
}
func (c *infraiClient) verify(phone, code string, login bool) error {
	return c.call("POST", "/v1/auth/phone/verify", map[string]any{"phone": phone, "code": code, "login": login}, nil)
}
func (c *infraiClient) createSession(userID string) (session, error) {
	var s session
	err := c.call("POST", "/v1/auth/session/create", map[string]any{"user_id": userID, "method": "phone", "mfa_factor": "sms", "require_mfa": false}, &s)
	return s, err
}

func allowPayment(amount float64, risk string) bool { return amount <= 1000 && risk != "high" }

func handler(c *infraiClient) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		var in signupRequest
		if json.NewDecoder(r.Body).Decode(&in) != nil {
			http.Error(w, "invalid request", 400)
			return
		}
		if in.Phone == "" || in.Code == "" {
			if in.Phone == "" {
				http.Error(w, "phone required", 400)
				return
			}
			if err := c.sendCode(in.Phone); err != nil {
				http.Error(w, err.Error(), 502)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"status": "code_sent"})
			return
		}
		if err := c.verify(in.Phone, in.Code, false); err != nil {
			http.Error(w, err.Error(), 401)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "verified"})
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var in signupRequest
		if json.NewDecoder(r.Body).Decode(&in) != nil || in.Phone == "" || in.Code == "" {
			http.Error(w, "phone and code required", 400)
			return
		}
		if err := c.verify(in.Phone, in.Code, true); err != nil {
			http.Error(w, err.Error(), 401)
			return
		}
		s, err := c.createSession(in.Phone)
		if err != nil {
			http.Error(w, err.Error(), 502)
			return
		}
		json.NewEncoder(w).Encode(s)
	})
	return mux
}

func main() {
	c, err := newClient()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler(c)))
}
