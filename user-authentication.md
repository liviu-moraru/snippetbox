# Chapter 11. User authentication

## 11.2 Creating a user model

```mysql
CREATE TABLE users (
id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
name VARCHAR(255) NOT NULL,
email VARCHAR(255) NOT NULL,
hashed_password CHAR(60) NOT NULL,
created DATETIME NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
```
### 11.3 User signup and password encryption

- We're not re-displaying the password if the form fails validations. This is because we don't want there to be [any risk](https://ux.stackexchange.com/questions/20418/when-form-submission-fails-password-field-gets-blanked-why-is-that-the-case) of the browser caching the plain-text password entered by the user.
- Validator: Matches regular expression

```
validator.go

func (v *Validator) Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
```
- #### Email pattern recommended by the W#C and Web Hypertext Application Technology Working Group: [Here](https://html.spec.whatwg.org/multipage/input.html#valid-e-mail-address) or [here](https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html#email-state-typeemail)

```
validator.go

// EmailRX Use the regexp.MustCompile() function to parse a regular expression pattern
// for sanity checking the format of an email address. This returns a pointer to
// a 'compiled' regexp.Regexp type, or panics in the event of an error. Parsing
// this pattern once at startup and storing the compiled *regexp.Regexp in a
// variable is more performant than re-parsing the pattern each time we need it.
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

handlers.go

form.CheckField(form.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
```
- #### In the solution, password encryption was implemented using bcrypt algorithm in the golang.org.x.crypto package
    - Import the latest version of the golang.org/x/crypto/bcrypt package

```
go get golang.org/x/crypto/bcrypt@latest
```
    - Test encrypt
```
validator_test.go

func TestPasswordHashFunction(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("my plain text password"), 12)
	if err != nil {
		t.Fatal("Encrypt error")
	}
	fmt.Println(string(hash), len(hash))
	// Result: $2a$12$EKbiojA0pFriBB5miZCzOuJGwan3IrqDzqj3DF7jQKo2P5NHikRTa 60
	hash2 := []byte(string(hash))
	err = bcrypt.CompareHashAndPassword(hash2, []byte("my plain text password"))
	if err != nil {
		t.Fatal("Comparison error")
	}
}
	
```
  - [Bcrypt output format](https://blog.boot.dev/cryptography/bcrypt-step-by-step/):

$2a$12$EKbiojA0pFriBB5miZCzOuJGwan3IrqDzqj3DF7jQKo2P5NHikRTa 
String contains 3 sections separated by $:
- First section 2a - hash algorithm identifier (len 3 bytes including $)
- Second section 12 - Cost factor ( 2^12 rounds of key expansions) (len 3 bytes including $)
- Third section: salt + hash (len 22 + 31 + 1 = 54 bytes, including $)
    - EKbiojA0pFriBB5miZCzOu 16-byte salt, base64 encoded to 22 characters
    - JGwan3IrqDzqj3DF7jQKo2P5NHikRTa 24 byte hash, base64 encoded to 31 characters
Total length of output: 3 + 3 + 54 = 60 bytes

- A cost of 12 is a reasonable minimum, but if possible you should carry out load testing, and if you can set the cost higher without adversely affecting user experience then you should.


#### [MySQL's error codes](https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html) and how they are treated

- 

