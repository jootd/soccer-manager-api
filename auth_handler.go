package main

// type AuthHandler struct {
// 	userBus userbus.ExtBusiness
// 	teamBus *teambus.Business
// }

// func NewAuthHandler(userBus userbus.ExtBusiness, teamBus *teambus.Business) *AuthHandler {
// 	return &AuthHandler{
// 		userBus: userBus,
// 		teamBus: teamBus,
// 	}
// }

// var jwtKey = []byte("test_secret_key")

// func (ah *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "cannot read body", http.StatusInternalServerError)
// 		return
// 	}

// 	var req RegisterRequest
// 	if err := json.Unmarshal(body, &req); err != nil {
// 		http.Error(w, "invalid JSON", http.StatusBadRequest)
// 		return
// 	}

// 	if len(req.Username) == 0 || len(req.Password) == 0 {
// 		http.Error(w, "username and password required", http.StatusBadRequest)
// 		return
// 	}

// 	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		http.Error(w, "failed to hash password", http.StatusInternalServerError)
// 		return
// 	}

// 	err = ah.userBus.Create(r.Context(), userbus.CreateUser{
// 		Username:     req.Username,
// 		PasswordHash: string(hash),
// 	})
// 	if err != nil {
// 		http.Error(w, "user already exists", http.StatusConflict)
// 		return
// 	}
// 	team, err := ah.teamBus.Create(r.Context())
// 	if err != nil {
// 		http.Error(w, "something went wrong, please try again", http.StatusInternalServerError)
// 		return
// 	}

// 	var user userbus.User
// 	err = ah.userBus.Update(r.Context(), user, userbus.UpdateUser{
// 		Username: &req.Username,
// 		TeamID:   &team.ID,
// 	})
// 	if err != nil {
// 		http.Error(w, "something went wrong, please try again", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Write([]byte(fmt.Sprintf("registered, user table %+v", user)))
// }

// func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	body, err := io.ReadAll(r.Body)
// 	defer r.Body.Close()
// 	if err != nil {
// 		http.Error(w, "cannot read body", http.StatusInternalServerError)
// 		return
// 	}

// 	var req RegisterRequest
// 	if err := json.Unmarshal(body, &req); err != nil {
// 		http.Error(w, "invalid JSON", http.StatusBadRequest)
// 		return
// 	}

// 	user, err := ah.userBus.Get(r.Context(), req.Username)
// 	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
// 		http.Error(w, "username or password incorrect", http.StatusUnauthorized)
// 		return
// 	}

// 	token, err := generateJWT(req.Username)
// 	if err != nil {
// 		http.Error(w, "something went wrong, please try again", http.StatusInternalServerError)
// 		return
// 	}

// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "donottouchme",
// 		Value:    token,
// 		Path:     "/",
// 		HttpOnly: true,
// 		Secure:   true,
// 		SameSite: http.SameSiteStrictMode,
// 		Expires:  time.Now().Add(1 * time.Hour),
// 	})

// 	w.Write([]byte(fmt.Sprintf("Logged in. Token: %s, userStore %+v", token, user)))
// }
