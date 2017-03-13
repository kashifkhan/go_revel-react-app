package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	//"io/ioutil"
	"posts/app/models"
	"strconv"
	//"strings"
)

type PostController struct {
	GorpController
}

func (p PostController) parsePost() (models.Post, error) {
	post := models.Post{}
	err := json.NewDecoder(p.Request.Body).Decode(&post)
	revel.INFO.Println(post.Title)

	return post, err
}

func (p PostController) Add() revel.Result {
	if post, err := p.parsePost(); err != nil {
		return p.RenderText("Unable to parse the Post from JSON.")
	} else {
		// Validate the model
		post.Validate(p.Validation)
		if p.Validation.HasErrors() {
			// Do something better here!
			return p.RenderText("You have error with the User.")
		} else {
			if err := p.Txn.Insert(&post); err != nil {
				return p.RenderText(
					"Error inserting record into database!")
			} else {
				return p.RenderJson(post)
			}
		}
	}
}

func (p PostController) Get(id int64) revel.Result {
	post := new(models.Post)
	err := p.Txn.SelectOne(post, `SELECT * FROM Post WHERE Id = ?`, id)
	if err != nil {
		return p.RenderText("Error.  Post probably doesn't exist.")
	}
	return p.RenderJson(post)
}

func (p PostController) List() revel.Result {
	lastId := parseIntOrDefault(p.Params.Get("lid"), -1)
	limit := parseUintOrDefault(p.Params.Get("limit"), uint64(25))
	posts, err := p.Txn.Select(models.Post{},
		`SELECT * FROM Post WHERE Id > ? LIMIT ?`, lastId, limit)
	if err != nil {
		return p.RenderText("Error trying to get records from DB.")
	}
	return p.RenderJson(posts)
}

func (p PostController) Update(id int64) revel.Result {
	post, err := p.parsePost()
	if err != nil {
		return p.RenderText("Unable to parse the User from JSON.")
	}
	// Ensure the Id is set.
	post.Id = id
	success, err := p.Txn.Update(&post)
	if err != nil || success == 0 {
		return p.RenderText("Unable to update user.")
	}
	return p.RenderText("Updated %v", id)
}

func (p PostController) Delete(id int64) revel.Result {
	success, err := p.Txn.Delete(&models.Post{Id: id})
	if err != nil || success == 0 {
		return p.RenderText("Failed to remove User")
	}
	return p.RenderText("Deleted %v", id)
}

func parseUintOrDefault(intStr string, _default uint64) uint64 {
	if value, err := strconv.ParseUint(intStr, 0, 64); err != nil {
		return _default
	} else {
		return value
	}
}

func parseIntOrDefault(intStr string, _default int64) int64 {
	if value, err := strconv.ParseInt(intStr, 0, 64); err != nil {
		return _default
	} else {
		return value
	}
}
