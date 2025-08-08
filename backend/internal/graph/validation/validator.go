package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"backend/internal/graph/errors"
	"backend/internal/graph/model"
	"backend/internal/graph/scalars"
)

// Validator provides input validation for GraphQL operations
type Validator struct{}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateCreatePostInput validates post creation input
func (v *Validator) ValidateCreatePostInput(input model.CreatePostInput) error {
	if err := v.ValidateTitle(input.Title); err != nil {
		return err
	}
	
	if err := v.ValidateContent(input.Content); err != nil {
		return err
	}
	
	if err := v.ValidateTags(input.Tags); err != nil {
		return err
	}
	
	return nil
}

// ValidateUpdatePostInput validates post update input
func (v *Validator) ValidateUpdatePostInput(input model.UpdatePostInput) error {
	if input.Title != nil {
		if err := v.ValidateTitle(*input.Title); err != nil {
			return err
		}
	}
	
	if input.Content != nil {
		if err := v.ValidateContent(*input.Content); err != nil {
			return err
		}
	}
	
	if input.Tags != nil {
		if err := v.ValidateTags(input.Tags); err != nil {
			return err
		}
	}
	
	return nil
}

// ValidateCreateUserInput validates user creation input
func (v *Validator) ValidateCreateUserInput(input model.CreateUserInput) error {
	if err := v.ValidateEmail(input.Email); err != nil {
		return err
	}
	
	if err := v.ValidateName(input.Name); err != nil {
		return err
	}
	
	if err := v.ValidatePassword(input.Password); err != nil {
		return err
	}
	
	return nil
}

// ValidateUpdateUserInput validates user update input
func (v *Validator) ValidateUpdateUserInput(input model.UpdateUserInput) error {
	if input.Name != nil {
		if err := v.ValidateName(*input.Name); err != nil {
			return err
		}
	}
	
	if input.Avatar != nil {
		if err := v.ValidateAvatarURL(*input.Avatar); err != nil {
			return err
		}
	}
	
	return nil
}

// ValidateTitle validates post title
func (v *Validator) ValidateTitle(title string) error {
	title = strings.TrimSpace(title)
	
	if title == "" {
		return errors.NewValidationError("Title cannot be empty", "title")
	}
	
	if utf8.RuneCountInString(title) < 3 {
		return errors.NewValidationError("Title must be at least 3 characters long", "title")
	}
	
	if utf8.RuneCountInString(title) > 200 {
		return errors.NewValidationError("Title cannot exceed 200 characters", "title")
	}
	
	// Check for invalid characters
	if strings.ContainsAny(title, "<>\"'&") {
		return errors.NewValidationError("Title contains invalid characters", "title")
	}
	
	return nil
}

// ValidateContent validates post content
func (v *Validator) ValidateContent(content string) error {
	content = strings.TrimSpace(content)
	
	if content == "" {
		return errors.NewValidationError("Content cannot be empty", "content")
	}
	
	if utf8.RuneCountInString(content) < 10 {
		return errors.NewValidationError("Content must be at least 10 characters long", "content")
	}
	
	if utf8.RuneCountInString(content) > 50000 {
		return errors.NewValidationError("Content cannot exceed 50,000 characters", "content")
	}
	
	return nil
}

// ValidateTags validates post tags
func (v *Validator) ValidateTags(tags []string) error {
	if len(tags) > 10 {
		return errors.NewValidationError("Cannot have more than 10 tags", "tags")
	}
	
	tagMap := make(map[string]bool)
	tagRegex := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	
	for i, tag := range tags {
		tag = strings.TrimSpace(strings.ToLower(tag))
		
		if tag == "" {
			return errors.NewValidationError(fmt.Sprintf("Tag %d cannot be empty", i+1), "tags")
		}
		
		if utf8.RuneCountInString(tag) < 2 {
			return errors.NewValidationError(fmt.Sprintf("Tag '%s' must be at least 2 characters long", tag), "tags")
		}
		
		if utf8.RuneCountInString(tag) > 30 {
			return errors.NewValidationError(fmt.Sprintf("Tag '%s' cannot exceed 30 characters", tag), "tags")
		}
		
		if !tagRegex.MatchString(tag) {
			return errors.NewValidationError(fmt.Sprintf("Tag '%s' contains invalid characters (only letters, numbers, hyphens, and underscores allowed)", tag), "tags")
		}
		
		if tagMap[tag] {
			return errors.NewValidationError(fmt.Sprintf("Duplicate tag '%s'", tag), "tags")
		}
		
		tagMap[tag] = true
	}
	
	return nil
}

// ValidateEmail validates email address
func (v *Validator) ValidateEmail(email string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	
	if err := scalars.ValidateEmail(email); err != nil {
		return errors.NewValidationError(err.Error(), "email")
	}
	
	return nil
}

// ValidateName validates user name
func (v *Validator) ValidateName(name string) error {
	name = strings.TrimSpace(name)
	
	if name == "" {
		return errors.NewValidationError("Name cannot be empty", "name")
	}
	
	if utf8.RuneCountInString(name) < 2 {
		return errors.NewValidationError("Name must be at least 2 characters long", "name")
	}
	
	if utf8.RuneCountInString(name) > 100 {
		return errors.NewValidationError("Name cannot exceed 100 characters", "name")
	}
	
	// Check for invalid characters (allow letters, spaces, hyphens, apostrophes)
	nameRegex := regexp.MustCompile(`^[a-zA-Z\s\-'\.]+$`)
	if !nameRegex.MatchString(name) {
		return errors.NewValidationError("Name contains invalid characters", "name")
	}
	
	return nil
}

// ValidatePassword validates password strength
func (v *Validator) ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.NewValidationError("Password must be at least 8 characters long", "password")
	}
	
	if len(password) > 128 {
		return errors.NewValidationError("Password cannot exceed 128 characters", "password")
	}
	
	// Check for at least one letter and one number
	hasLetter := false
	hasNumber := false
	
	for _, char := range password {
		if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' {
			hasLetter = true
		} else if char >= '0' && char <= '9' {
			hasNumber = true
		}
	}
	
	if !hasLetter {
		return errors.NewValidationError("Password must contain at least one letter", "password")
	}
	
	if !hasNumber {
		return errors.NewValidationError("Password must contain at least one number", "password")
	}
	
	// Optional: require special character for stronger passwords
	// if !hasSpecial {
	//     return errors.NewValidationError("Password must contain at least one special character", "password")
	// }
	
	return nil
}

// ValidateAvatarURL validates avatar URL
func (v *Validator) ValidateAvatarURL(url string) error {
	url = strings.TrimSpace(url)
	
	if url == "" {
		return nil // Avatar is optional
	}
	
	if len(url) > 2048 {
		return errors.NewValidationError("Avatar URL cannot exceed 2048 characters", "avatar")
	}
	
	// Simple URL validation
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	if !urlRegex.MatchString(url) {
		return errors.NewValidationError("Invalid avatar URL format", "avatar")
	}
	
	// Check for image file extensions
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg"}
	hasValidExtension := false
	lowerURL := strings.ToLower(url)
	
	for _, ext := range imageExtensions {
		if strings.Contains(lowerURL, ext) {
			hasValidExtension = true
			break
		}
	}
	
	if !hasValidExtension {
		return errors.NewValidationError("Avatar URL must point to an image file", "avatar")
	}
	
	return nil
}

// ValidateCommentContent validates comment content
func (v *Validator) ValidateCommentContent(content string) error {
	content = strings.TrimSpace(content)
	
	if content == "" {
		return errors.NewValidationError("Comment content cannot be empty", "content")
	}
	
	if utf8.RuneCountInString(content) < 1 {
		return errors.NewValidationError("Comment must have at least 1 character", "content")
	}
	
	if utf8.RuneCountInString(content) > 2000 {
		return errors.NewValidationError("Comment cannot exceed 2000 characters", "content")
	}
	
	return nil
}

// ValidateSearchQuery validates search query
func (v *Validator) ValidateSearchQuery(query string) error {
	query = strings.TrimSpace(query)
	
	if query == "" {
		return errors.NewValidationError("Search query cannot be empty", "query")
	}
	
	if utf8.RuneCountInString(query) < 2 {
		return errors.NewValidationError("Search query must be at least 2 characters long", "query")
	}
	
	if utf8.RuneCountInString(query) > 100 {
		return errors.NewValidationError("Search query cannot exceed 100 characters", "query")
	}
	
	return nil
}

// ValidatePaginationInput validates pagination parameters
func (v *Validator) ValidatePaginationInput(input *model.PaginationInput) error {
	if input == nil {
		return nil
	}
	
	if input.Page != nil && *input.Page < 1 {
		return errors.NewValidationError("Page must be at least 1", "page")
	}
	
	if input.Page != nil && *input.Page > 1000 {
		return errors.NewValidationError("Page cannot exceed 1000", "page")
	}
	
	if input.Limit != nil && *input.Limit < 1 {
		return errors.NewValidationError("Limit must be at least 1", "limit")
	}
	
	if input.Limit != nil && *input.Limit > 100 {
		return errors.NewValidationError("Limit cannot exceed 100", "limit")
	}
	
	return nil
}