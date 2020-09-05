package main

import (
	"strconv"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
)

// Todo struct for all todos
type Todo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

var todos = []*Todo{
	{ID: 1, Name: "Walk the dog", Completed: false},
	{ID: 2, Name: "Walk the cat", Completed: false},
	{ID: 3, Name: "Walk the bat", Completed: true},
}

func main() {
	app := fiber.New()

	app.Use(middleware.Logger())

	app.Get("/", func(ctx *fiber.Ctx) {
		ctx.Send("hello world")
	})

	setupAPIV1(app)

	app.Listen(3000)

}

func setupAPIV1(app *fiber.App) {
	v1 := app.Group("/v1")
	setupTodosRoutes(v1)
}

func setupTodosRoutes(grp fiber.Router) {
	todosRoutes := grp.Group("/todos")
	todosRoutes.Get("/", GetTodos)
	todosRoutes.Get("/:id", GetTodo)
	todosRoutes.Post("/", CreateTodo)
	todosRoutes.Delete("/:id", DeleteTodo)
	todosRoutes.Patch("/:id", UpdateTodo)
}

// UpdateTodo will update the information to todo
func UpdateTodo(ctx *fiber.Ctx) {

	type request struct {
		Name      *string `json:"name"`
		Completed *bool   `json:"completed`
	}

	paramID := ctx.Params("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON((fiber.Map{
			"error": "cannot parse id",
		}))
		return
	}

	var body request
	err = ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse body",
		})
		return
	}

	var todo *Todo

	for _, t := range todos {
		if t.ID == id {
			todo = t
			break
		}
	}

	if todo == nil {
		ctx.Status(fiber.StatusNotFound)
		return
	}

	if body.Name != nil {
		todo.Name = *body.Name
	}

	if body.Completed != nil {
		todo.Completed = *body.Completed
	}

	// you surely think the todo is update but....

	ctx.Status(fiber.StatusOK).JSON(todo)
}

// DeleteTodo will delete the todos by id
func DeleteTodo(ctx *fiber.Ctx) {
	paramID := ctx.Params("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[0:i], todos[i+1:]...)
			ctx.Status(fiber.StatusNoContent)
			return
		}
	}

	ctx.Status(fiber.StatusNotFound)

}

// CreateTodo will crate the todo to the list
func CreateTodo(ctx *fiber.Ctx) {
	type request struct {
		Name string `json:"name"`
	}

	var body request

	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
		})
		return
	}

	todo := &Todo{
		ID:        len(todos) + 1,
		Name:      body.Name,
		Completed: false,
	}

	todos = append(todos, todo)

	ctx.Status(fiber.StatusCreated).JSON(todos)

}

// GetTodo will get only one record from id
func GetTodo(ctx *fiber.Ctx) {
	paramID := ctx.Params("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return
	}

	for _, todo := range todos {
		if todo.ID == id {
			ctx.Status(fiber.StatusOK).JSON(todo)
			return
		}
	}

	ctx.Status(fiber.StatusNotFound)
}

// GetTodos will get all todos
func GetTodos(ctx *fiber.Ctx) {
	ctx.Status(fiber.StatusOK).JSON(todos)
}
