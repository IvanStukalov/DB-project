package main

import (
	"context"
	forumDelivery "github.com/IvanStukalov/DB_project/internal/pkg/forum/delivery"
	forumRepo "github.com/IvanStukalov/DB_project/internal/pkg/forum/repo"
	forumUsecase "github.com/IvanStukalov/DB_project/internal/pkg/forum/usecase"
	postDelivery "github.com/IvanStukalov/DB_project/internal/pkg/post/delivery"
	postRepo "github.com/IvanStukalov/DB_project/internal/pkg/post/repo"
	postUsecase "github.com/IvanStukalov/DB_project/internal/pkg/post/usecase"
	serviceDelivery "github.com/IvanStukalov/DB_project/internal/pkg/service/delivery"
	serviceRepo "github.com/IvanStukalov/DB_project/internal/pkg/service/repo"
	serviceUsecase "github.com/IvanStukalov/DB_project/internal/pkg/service/usecase"
	threadDelivery "github.com/IvanStukalov/DB_project/internal/pkg/thread/delivery"
	threadRepo "github.com/IvanStukalov/DB_project/internal/pkg/thread/repo"
	threadUsecase "github.com/IvanStukalov/DB_project/internal/pkg/thread/usecase"
	userDelivery "github.com/IvanStukalov/DB_project/internal/pkg/user/delivery"
	userRepo "github.com/IvanStukalov/DB_project/internal/pkg/user/repo"
	userUsecase "github.com/IvanStukalov/DB_project/internal/pkg/user/usecase"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
)

// sudo docker rm -f my_container
// sudo docker build -t docker .
// sudo docker run -p 5000:5000 --name my_container -t docker
// ./technopark-dbms-forum func -u http://localhost:5000/api -r report.html

func main() {
	muxRoute := mux.NewRouter()
	//conn := "postgres://postgres:password@127.0.0.1:5432/bd?sslmode=disable&pool_max_conns=1000"
	conn := "postgres://docker:docker@127.0.0.1:5432/docker?sslmode=disable&pool_max_conns=1000"
	pool, err := pgxpool.Connect(context.Background(), conn)
	if err != nil {
		log.Fatal("No connection to postgres", err)
	}

	uRepo := userRepo.NewRepoPostgres(pool)
	tRepo := threadRepo.NewRepoPostgres(pool)
	fRepo := forumRepo.NewRepoPostgres(pool)
	pRepo := postRepo.NewRepoPostgres(pool)
	sRepo := serviceRepo.NewRepoPostgres(pool)

	uUsecase := userUsecase.NewRepoUsecase(uRepo)
	tUsecase := threadUsecase.NewRepoUsecase(tRepo, pRepo, uRepo)
	fUsecase := forumUsecase.NewRepoUsecase(fRepo, uRepo, tRepo)
	pUsecase := postUsecase.NewRepoUsecase(pRepo, fRepo, uRepo, tRepo)
	sUsecase := serviceUsecase.NewRepoUsecase(sRepo)

	uHandler := userDelivery.NewUserHandler(uUsecase)
	tHandler := threadDelivery.NewThreadHandler(tUsecase)
	fHandler := forumDelivery.NewForumHandler(fUsecase)
	pHandler := postDelivery.NewPostHandler(pUsecase)
	sHandler := serviceDelivery.NewServiceHandler(sUsecase)

	forum := muxRoute.PathPrefix("/api").Subrouter()
	{
		forum.HandleFunc("/user/{nickname}/create", uHandler.CreateUser).Methods(http.MethodPost)
		forum.HandleFunc("/user/{nickname}/profile", uHandler.GetUser).Methods(http.MethodGet)
		forum.HandleFunc("/user/{nickname}/profile", uHandler.UpdateUser).Methods(http.MethodPost)

		forum.HandleFunc("/forum/create", fHandler.CreateForum).Methods(http.MethodPost)
		forum.HandleFunc("/forum/{slug}/details", fHandler.GetForum).Methods(http.MethodGet)
		forum.HandleFunc("/forum/{slug}/threads", fHandler.GetForumThreads).Methods(http.MethodGet)
		forum.HandleFunc("/forum/{slug}/create", fHandler.CreateThread).Methods(http.MethodPost)
		forum.HandleFunc("/forum/{slug}/users", fHandler.GetUsers).Methods(http.MethodGet)

		forum.HandleFunc("/thread/{slug_or_id}/details", tHandler.GetThread).Methods(http.MethodGet)
		forum.HandleFunc("/thread/{slug_or_id}/details", tHandler.UpdateThread).Methods(http.MethodPost)
		forum.HandleFunc("/thread/{slug_or_id}/create", tHandler.CreatePosts).Methods(http.MethodPost)
		forum.HandleFunc("/thread/{slug_or_id}/vote", tHandler.CreateVote).Methods(http.MethodPost)
		forum.HandleFunc("/thread/{slug_or_id}/posts", tHandler.GetPosts).Methods(http.MethodGet)

		forum.HandleFunc("/post/{id}/details", pHandler.GetPost).Methods(http.MethodGet)
		forum.HandleFunc("/post/{id}/details", pHandler.UpdatePost).Methods(http.MethodPost)

		forum.HandleFunc("/service/status", sHandler.Status).Methods(http.MethodGet)
		forum.HandleFunc("/service/clear", sHandler.Clear).Methods(http.MethodPost)
	}

	http.Handle("/", muxRoute)
	log.Print(http.ListenAndServe(":5000", muxRoute))
}
