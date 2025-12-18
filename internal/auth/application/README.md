# Application

Здесь расположен Application-слой по DDD. 
Этот слой реализует пользовательские сценарии, то есть то, что мы называем use case.
По сути, мы просто оркестрируем действия domain слоя.

Вместо того, чтобы заводить один огромный интерфейс для этого слоя, типа такого

``` Go
type Application interface {
    // аргументы и возвр. значения не указаны для удобства
    Login()
    Logout()
    // другие методы
}
```

используются интерфейсы такого вида

``` Go
type Logout interface {
    Execute(ctx context.Context, cmd *LogoutCommand) (*LogoutRespone, error)
}

// передаваемые данные
type LogoutCommand struct {
    // поля структуры
}

// возвращаемые данные
type LogoutResponse struct {
    // поля структуры
} 
```

Это сделано для того, чтобы не делать огромный интерфейс, который бы нарушал SOLID.
А именно принцип Interface Segregation Principle (возможно ещё и SRP, но это не уверен)

Кроме того, так удобнее тестировать.