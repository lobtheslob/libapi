CREATE TABLE `book` (
	  `id` varchar(36),
		`name` text,
		`author` text,
		`isbn` varchar(36),
		PRIMARY KEY `book_primary_key` (`id`)
				) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
