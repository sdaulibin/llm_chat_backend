-- 创建文件信息表
CREATE TABLE files (
	id serial NOT NULL,
    user_id varchar(10),
	file_id uuid NOT NULL,
	original_filename varchar(256) NOT NULL,
	local_filename varchar(256) NOT NULL,
	file_path varchar(512) NOT NULL,
	file_size bigint NOT NULL,
	file_type varchar(50) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP(0) NOT NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP(0) NOT NULL,
	CONSTRAINT file_pkey PRIMARY KEY (id)
);