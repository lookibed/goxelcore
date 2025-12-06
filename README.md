

## Цель: Портирование C++ кода из [`./voxelcore/`](https://github.com/MihailRis/voxelcore) в `./goxelcore/`.

### Принципы:
- Структура файлов и названия сохраняются 1:1.
- Используется `go-sdl3` (уже настроен).
- Для звука и шрифтов использовать заглушки.

---

## Выполнено

- **Структура проекта:** Создана зеркальная иерархия папок и пустые Go-заглушки для всех C++ файлов.
- **Настройки и База:** Портированы `CoreParameters`, `Settings`, `EngineSettings`, `Typedefs`.
- **Окно и Ввод:** Реализованы `Window`, `WindowControl` и подсистема `Input` (SDL3).
- **Движок (Engine):** Настроена инициализация контекста OpenGL, тайминги (`Time`), логирование (`Logger`) и пути к ресурсам (`ResPaths`).
- **Графическое ядро:** Портированы компоненты:
    - `Shader` (включая парсинг GLSL), `DrawContext`, `Commons`.
    - `ImageData`, `Texture`, `UVRegion`.
    - `Mesh`, `MeshData`.
    - `Batch2D`, `Batch3D`, `Framebuffer`.
- **IO:** Реализованы работа с путями и файлами (`io/path`, `io/io`).

---

## Следующие шаги

1.  **Портировать `EngineController`** (`logic/EngineController`): Основная логика управления движком.
2.  **Портировать `Assets`** (`assets/Assets`): Загрузчики ресурсов.
3.  **Портировать `ContentControl`** (`content/ContentControl`): Управление контентом игры.
4.  **Портировать `GUI`** (`graphics/ui/GUI`): Пользовательский интерфейс.