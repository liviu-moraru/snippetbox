package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/liviu-moraru/snippetbox/internal/models"
	"github.com/liviu-moraru/snippetbox/internal/validator"
	"net/http"
	"os"
	"path"
	"strconv"
)

func (app *Application) HomeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		snippets, err := app.Snippets.Latest()
		if err != nil {
			app.serverError(w, err)
			return
		}

		data := app.newTemplateData(r)
		data.Snippets = snippets

		app.render(w, http.StatusOK, "home.tmpl", data)
	})
}

func (app *Application) SnippetViewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())

		id, err := strconv.Atoi(params.ByName("id"))

		if err != nil || id < 1 {
			app.notFound(w)
			return
		}

		snippet, err := app.Snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}
			return
		}

		data := app.newTemplateData(r)
		data.Snippet = snippet

		app.render(w, http.StatusOK, "view.tmpl", data)
	})
}

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl", data)
}

// Update our snippetCreateForm struct to include struct tags which tell the
// decoder how to map HTML form values into the different struct fields. So, for
// example, here we're telling the decoder to store the value from the HTML form
// input with the name "title" in the Title field. The struct tag `form:"-"`
// tells the decoder to completely ignore a field during decoding.
type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires""`
	validator.Validator `form:"-"`
}

func (app *Application) SnippetCreatePostHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var form snippetCreateForm

		err := app.decodePostForm(r, &form)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		// Then validate and use the data as normal...
		form.CheckField(form.NotBlank(form.Title), "title", "This field cannot be blank")
		form.CheckField(form.MaxCharacters(form.Title, 100), "title", "This field cannot be more than 100 characters long")
		form.CheckField(form.NotBlank(form.Content), "content", "This field cannot be blank")
		form.CheckField(form.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

		// If there are any validation errors re-display the create.tmpl template,
		// passing in the snippetCreateForm instance as dynamic data in the Form
		// field. Note that we use the HTTP status code 422 Unprocessable Entity
		// when sending the response to indicate that there was a validation error.
		if !form.Valid() {
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
			return
		}

		// We also need to update this line to pass the data from the
		// snippetCreateForm instance to our Insert() method.
		id, err := app.Snippets.Insert(form.Title, form.Content, form.Expires)
		if err != nil {
			app.serverError(w, err)
			return
		}

		// Use the Put() method to add a string value ("Snippet successfully
		// created!") and the corresponding key ("flash") to the session data.
		app.SessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	})
}

// Create a new userSignupForm struct.
type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// Update the handler so it displays the signup page.
func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl", data)
}

func (app *Application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	// Parse the form data into the userSignupForm struct.
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using our helper functions.
	form.CheckField(form.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(form.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(form.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(form.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(form.MinChars(form.Password, 8), "password", "This field must be at least 8 character long")

	// If there are any errors, redisplay the signup form along with a 422
	// status code.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	// Try to create a new user record in the database. If the email already
	// exists then add an error message to the form and re-display it.
	err = app.Users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email is already is use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, err)
		}

		return
	}

	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked.
	app.SessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

func (app *Application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for logging in a user...")
}

func (app *Application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *Application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}

func (app *Application) NoDirListingHandler(d http.Dir) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		fp := params.ByName("filepath")
		fp = path.Join(string(d), fp)
		if fi, err := os.Stat(fp); err == nil && !fi.IsDir() {
			fileServer := http.FileServer(d)
			http.StripPrefix("/static", fileServer).ServeHTTP(w, r)
			return
		}
		app.notFound(w)
	})
}

func httpRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r,
		"https://"+r.Host+r.URL.String(),
		http.StatusMovedPermanently)
}
