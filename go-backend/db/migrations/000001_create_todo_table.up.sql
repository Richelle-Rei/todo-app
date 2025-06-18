CREATE TABLE IF NOT EXISTS todo(
id BIGSERIAL PRIMARY KEY NOT NULL,
user_id BIGINT NOT NULL,                            
title TEXT NOT NULL,                                                                                             
completed BOOLEAN NOT NULL,
display_order BIGINT UNIQUE NOT NULL 
);

INSERT INTO todo (user_id, title, completed, display_order) VALUES 
(1, 'delectus aut autem', false, 1),
(1, 'quis ut nam facilis et officia qui', false, 2),
(1, 'fugiat veniam minus', false, 3),
(1, 'et porro tempora', true, 4),
(1, 'laboriosam mollitia et enim quasi adipisci quia provident illum', false, 5);