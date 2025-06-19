CREATE TABLE IF NOT EXISTS todo(
id BIGSERIAL PRIMARY KEY NOT NULL,
user_id BIGINT NOT NULL,                            
title TEXT NOT NULL,                                                                                             
completed BOOLEAN NOT NULL,
display_order BIGINT UNIQUE NOT NULL 
);