-- Skema awal. HARUS konsisten dengan struct di internal/models.

CREATE TABLE IF NOT EXISTS users (
  id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name          VARCHAR(100)    NOT NULL,
  email         VARCHAR(191)    NOT NULL,
  password_hash VARCHAR(255)    NOT NULL,
  role          VARCHAR(20)     NOT NULL DEFAULT 'STUDENT',
  created_at    DATETIME(3)     NULL,
  updated_at    DATETIME(3)     NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_users_email (email),
  CONSTRAINT chk_users_role CHECK (role IN ('STUDENT', 'ADMIN'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS courses (
  id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  title       VARCHAR(150)    NOT NULL,
  description TEXT            NULL,
  created_at  DATETIME(3)     NULL,
  updated_at  DATETIME(3)     NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS quizzes (
  id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  course_id   BIGINT UNSIGNED NOT NULL,
  title       VARCHAR(150)    NOT NULL,
  description TEXT            NOT NULL,
  `constraints` TEXT          NULL,
  created_at  DATETIME(3)     NULL,
  updated_at  DATETIME(3)     NULL,
  PRIMARY KEY (id),
  KEY idx_quizzes_course_id (course_id),
  CONSTRAINT fk_quizzes_course FOREIGN KEY (course_id) REFERENCES courses (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS test_cases (
  id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  quiz_id         BIGINT UNSIGNED NOT NULL,
  stdin           TEXT            NULL,
  expected_output TEXT            NOT NULL,
  is_hidden       TINYINT(1)      NOT NULL DEFAULT 1,
  source          VARCHAR(20)     NOT NULL DEFAULT 'MANUAL',
  created_at      DATETIME(3)     NULL,
  updated_at      DATETIME(3)     NULL,
  PRIMARY KEY (id),
  KEY idx_test_cases_quiz_id (quiz_id),
  CONSTRAINT fk_test_cases_quiz FOREIGN KEY (quiz_id) REFERENCES quizzes (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS submissions (
  id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id       BIGINT UNSIGNED NOT NULL,
  quiz_id       BIGINT UNSIGNED NOT NULL,
  source_code   MEDIUMTEXT      NOT NULL,
  status        VARCHAR(20)     NOT NULL,
  score         INT             NOT NULL DEFAULT 0,
  passed_count  INT             NOT NULL DEFAULT 0,
  total_count   INT             NOT NULL DEFAULT 0,
  error_message TEXT            NULL,
  created_at    DATETIME(3)     NULL,
  PRIMARY KEY (id),
  KEY idx_submissions_user_quiz (user_id, quiz_id),
  KEY idx_submissions_quiz_id (quiz_id),
  CONSTRAINT fk_submissions_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT fk_submissions_quiz FOREIGN KEY (quiz_id) REFERENCES quizzes (id) ON DELETE CASCADE,
  CONSTRAINT chk_submissions_status CHECK (status IN ('PASSED', 'WRONG_ANSWER', 'COMPILE_ERROR', 'RUNTIME_ERROR'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Fix user permissions for external / docker network access
CREATE USER IF NOT EXISTS 'autograder'@'%' IDENTIFIED BY 'autograder_pass';
GRANT ALL PRIVILEGES ON autograder_db.* TO 'autograder'@'%';
FLUSH PRIVILEGES;
