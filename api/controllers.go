package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"roommate-finder/db/repo"

	"github.com/gin-gonic/gin"
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

func (h *UserHandler) handleUserMatch(c *gin.Context) {
	var req repo.CreateMatchParams
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//score := CalculateMatchScore(*req.User1ID, *req.User2ID)

	resp := repo.CreateMatchParams{
		User1ID:    req.User1ID,
		User2ID:    req.User2ID,
		MatchScore: req.MatchScore,
	}

	match, err := h.querier.CreateMatch(c, resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, match)

}

// hashing password
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
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

// func CalculateMatchScore(user1, user2 string) int32 {
// 	score := int32(0)
// 	for key, value := range user1.Preferences {
// 		if user2.Preferences[key] == value {
// 			score += 10 // Increase match score for similar preferences
// 		}
// 	}
// 	return score
// }

// func CalculateScore(user1PrefsJSON, user2PrefsJSON string) (int, error) {

//     var user1Prefs, user2Prefs Preferences

//     // Parse JSONB strings
//     err := json.Unmarshal([]byte(user1PrefsJSON), &user1Prefs)
//     if err != nil {
//         return 0, err
//     }
//     err = json.Unmarshal([]byte(user2PrefsJSON), &user2Prefs)
//     if err != nil {
//         return 0, err
//     }

//     // Compute score based on preference overlap
//     score := 0
//     for key, weight1 := range user1Prefs {
//         if weight2, exists := user2Prefs[key]; exists {
//             score += int(weight1 * weight2 * 100) // Adjust scale as needed
//         }
//     }

//     return score, nil
// }
