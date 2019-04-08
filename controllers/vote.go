package controllers

import (
	"github.com/Betchika99/tp_db_project/models"
	"log"
)

func VoteCreate(vote models.Vote, slugOrId string) (models.Thread, error){
	log.Println("Create Vote")

	thread := models.Thread{}
	return thread, nil
	//if thread, err := strconv.Atoi(slugOrId); err != nil {
	//	threadModel, err := database.SelectThreadBySlug(slugOrId)
	//	if err != nil {
	//		return thread, err
	//	}
	//	thread = threadModel.Id
	//	forum = threadModel.Forum
	//} else {
	//	forumModel, err := SelectThreadById(thread)
	//	if err != nil {
	//		return nil, err
	//	}
	//	forum = forumModel.Forum
	//}
	//if thread == 0 {
	//	return nil, fmt.Errorf("There are not any thread with this slug or id")
	//}
	//
	//
	//forum, err := database.SelectForumBySlug(forumSlug)
	//if err != nil {
	//	log.Println("ERROR is", err.Error())
	//
	//	message := forumBySlugNotFound + forumSlug
	//
	//	jsonBody, err := json.Marshal(models.ModelError{message})
	//	if err != nil {
	//		log.Println(" json marshal fail")
	//		return
	//	}
	//
	//	ctx.SetContentType("application/json")
	//	ctx.SetStatusCode(fasthttp.StatusNotFound)
	//	ctx.Response.SetBody(jsonBody)
	//	return
	//}
	//thread.Forum = forum.Slug
	//
	//_, err = database.SelectOneUser(thread.Author)
	//if err != nil {
	//	log.Println("ERROR is", err.Error())
	//
	//	message := authorNotFound + thread.Author
	//
	//	jsonBody, err := json.Marshal(models.ModelError{message})
	//	if err != nil {
	//		log.Println(" json marshal fail")
	//		return
	//	}
	//
	//	ctx.SetContentType("application/json")
	//	ctx.SetStatusCode(fasthttp.StatusNotFound)
	//	ctx.Response.SetBody(jsonBody)
	//	return
	//}
	//
	//existThread := models.Thread{}
	//if thread.Slug != "" {
	//	existThread, err = database.SelectThreadBySlug(thread.Slug)
	//	if err != nil && err != pgx.ErrNoRows {
	//		log.Println("ERROR is", err.Error())
	//		return
	//	}
	//}
	//
	//if existThread.Id != 0 {
	//	jsonBody, err := json.Marshal(existThread)
	//	if err != nil {
	//		log.Println(" json marshal fail")
	//		return
	//	}
	//
	//	ctx.SetContentType("application/json")
	//	ctx.SetStatusCode(fasthttp.StatusConflict)
	//	ctx.Response.SetBody(jsonBody)
	//	return
	//}
	//
	//data, err := database.InsertThread(thread)
	//if err != nil {
	//	log.Println("ERROR is", err.Error())
	//	return
	//}
	//
	//jsonBody, err := json.Marshal(data)
	//if err != nil {
	//	log.Println(" json marshal fail")
	//	return
	//
	//}
	//ctx.SetContentType("application/json")
	//ctx.SetStatusCode(fasthttp.StatusCreated)
	//ctx.Response.SetBody(jsonBody)
}
