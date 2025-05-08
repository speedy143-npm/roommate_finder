package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"roommate-finder/db/repo"
	"roommate-finder/utils"
	"time"

	"github.com/gin-gonic/gin"
	//"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func (h *UserHandler) handleUserRegistration(c *gin.Context) {
	var req repo.CreateUserParams

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Validate number format
	if err := ValidateAndFormatNumber(req.Phoneno); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	HashedPassword, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate email format
	if err := ValidateEmail(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	preferencesJSON, err := json.Marshal(req.Preferences)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var preferences repo.PrefJson
	err = json.Unmarshal(preferencesJSON, &preferences)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	resp := repo.CreateUserParams{
		Fname:          req.Fname,
		Lname:          req.Lname,
		Phoneno:        req.Phoneno,
		Email:          req.Email,
		Password:       HashedPassword,
		Bio:            req.Bio,
		Preferences:    preferences,
		ProfilePicture: req.ProfilePicture,
	}

	user, err := h.querier.CreateUser(c, resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) handleGetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id is required"})
		return
	}

	user, errs := h.querier.GetUserById(c, id)
	if errs != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errs.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"User": user})

}

func (h *UserHandler) handleUpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id is required"})
		return
	}

	_, errs := h.querier.GetUserById(c, id)
	if errs != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "here"})
		return
	}

	var req repo.UpdateUserProfileParams
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, errs := h.querier.UpdateUserProfile(c, req)
	if errs != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errs.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"User Updated": user})
}

func (h *UserHandler) handleUserMatch(c *gin.Context) {
	id1 := c.Param("id1")
	id2 := c.Param("id2")

	if id1 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id1 is required"})
		return
	}

	if id2 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id2 is required"})
		return
	}

	user1, errs := h.querier.GetUserById(c, id1)
	if errs != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errs.Error()})
		return
	}

	user2, errs := h.querier.GetUserById(c, id2)
	if errs != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errs.Error()})
		return
	}

	var req repo.CreateMatchParams
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	score, cat, _ := CalculateScore(user1[0].Preferences, user2[0].Preferences)

	resp := repo.CreateMatchParams{
		User1ID:    &user1[0].ID,
		User2ID:    &user2[0].ID,
		MatchScore: score,
		Status:     cat,
	}

	match, err := h.querier.CreateMatch(c, resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"match": match})

}

func (h *UserHandler) handleForgotPassword(c *gin.Context) {
	type UserEmail struct {
		Email string `json:"email"`
	}

	var req UserEmail
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, errs := h.querier.GetUserByEmail(c, req.Email)
	if errs != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errs.Error()})
		return
	}

	// Generate a reset token
	resetToken, err := GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
		return
	}

	expiryTime := time.Now().Add(30 * time.Minute)

	// Convert to pgtype.Timestamp
	pgExpiryTime := pgtype.Timestamp{Time: expiryTime}
	fmt.Println(pgExpiryTime)

	// Store token in the database with expiration time
	resetParams := repo.ForgotPasswordParams{
		UserID: user.ID,
		Token:  resetToken,
		//Expiry: pgExpiryTime, // Token valid for 30 minutes
	}

	if _, err := h.querier.ForgotPassword(c, resetParams); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Construct the password reset link
	resetLink := fmt.Sprintf("https://roommatefinder.com/reset-password?token=%s", resetToken)
	fmt.Println(resetLink)

	// Send email with reset link (using a utility function)
	err = utils.SendEmail(user.Email, "Password Reset Request", fmt.Sprintf("Click here to reset your password: %s", resetLink))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) //"Failed to send reset email"
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}

func (h *UserHandler) handleResetPassword(c *gin.Context) {
	type ResetPasswordRequest struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Verify token exists and is not expired
	resetData, err := h.querier.GetResetToken(c, req.Token)
	if err != nil || resetData[0].Expiry.Time.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
		return
	}

	// Hash new password
	hashedPassword, err := HashPassword(req.NewPassword)

	var reqp = repo.UpdateUserPasswordParams{
		ID:       resetData[0].UserID,
		Password: hashedPassword,
	}

	// Update password in database
	_, err = h.querier.UpdateUserPassword(c, reqp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	// Delete reset token after successful password change
	err = h.querier.DeleteResetToken(c, req.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete reset token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password successfully reset"})
}

// hashing password
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "Failed to hash password", err
	}
	return string(hashedBytes), nil
}

// ComparePassword checks if the given password matches the stored hash
func ComparePassword(hashedPassword, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil // Returns true if passwords match
}

// ValidateEmail checks if the email is correctly formatted
func ValidateEmail(email string) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func ValidateAndFormatNumber(number string) error {
	re := regexp.MustCompile(`^(6(70|71|72|73|74|75|76|77|78|79|80|81|82|83|84|85|86|87|88|89)|65[0-9]|69[1-9]|62[0-3])\d{6}$`)

	if !re.MatchString(number) {
		return errors.New("invalid number format")
	}
	// Append country code if valid
	number = "237" + number
	return nil
}

// func SaveProfilePicture(userID string, imageURL string, db *sql.DB) error {
//     _, err := db.Exec("UPDATE users SET profile_picture = $1 WHERE id = $2", imageURL, userID)
//     return err
// }

func CalculateScore(user1Prefs, user2Prefs repo.PrefJson) (*int32, *string, error) {
	var score int32 = 0

	// Use reflection to iterate over struct fields dynamically
	user1Value := reflect.ValueOf(user1Prefs)
	user2Value := reflect.ValueOf(user2Prefs)
	user1Type := reflect.TypeOf(user1Prefs)

	totalPreferences := int32(user1Type.NumField()) // Count the number of preferences
	maxPossibleScore := totalPreferences * 10       // Maximum possible score (each preference contributes up to 10)

	for i := 0; i < user1Type.NumField(); i++ {
		user1Field := user1Value.Field(i).String()
		user2Field := user2Value.Field(i).String()

		// Apply scoring logic dynamically
		if user1Field == user2Field {
			score += 10 // Exact match
		} else if user1Field != "" && user2Field != "" {
			score += 5 // Both have a preference, but different values
		}
	}

	// Convert score to a percentage
	matchPercentage := (score * 100) / maxPossibleScore

	// Convert score to a percentage and format it as a string with a '%'
	// matchPercentage := fmt.Sprintf("%d%%", (score*100)/maxPossibleScore)

	// Determine category based on percentage
	var category string
	switch {
	case matchPercentage <= 30:
		category = "Poor"
	case matchPercentage > 30 && matchPercentage <= 60:
		category = "Good"
	case matchPercentage > 60 && matchPercentage <= 85:
		category = "Very Good"
	case matchPercentage > 85:
		category = "Excellent"
	}

	return &matchPercentage, &category, nil
}

// GenerateResetToken creates a secure random token
func GenerateToken() (string, error) {
	bytes := make([]byte, 32) // Generate 32 random bytes
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
