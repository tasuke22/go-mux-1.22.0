CREATE TABLE users
(
    id       CHAR(36) PRIMARY KEY,
    name     VARCHAR(255) NOT NULL,
    email    VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE todos
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    completed   BOOLEAN DEFAULT FALSE NOT NULL,
    user_id     CHAR(36)     NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX idx_user_id ON todos (user_id);


-- 1人のユーザーを作成
INSERT INTO users (id, name, email, password)
VALUES
    ('860c602c-2b99-48b3-82b7-253f52ad9c7e', 'tanaka', 'user@example.com', 'password123');

-- そのユーザーに紐づくタスクを2つ作成
INSERT INTO todos (title, description, completed, user_id)
VALUES
    ('タスク1のタイトル', 'タスク1の説明', FALSE, '860c602c-2b99-48b3-82b7-253f52ad9c7e'),
    ('タスク2のタイトル', 'タスク2の説明', FALSE, '860c602c-2b99-48b3-82b7-253f52ad9c7e');
