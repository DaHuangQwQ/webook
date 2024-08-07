-- 创建学校表
CREATE TABLE usr_schools (
                             school_id INT PRIMARY KEY AUTO_INCREMENT, 	-- 学校ID，主键，自动递增
                             school_name VARCHAR(100) NOT NULL, 					-- 学校名称，最多100个字符，不允许为空
                             address VARCHAR(255) NOT NULL 							-- 学校地址，最多255个字符，不允许为空
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = COMPACT;

-- 创建学生表
CREATE TABLE usr_students (
                              student_id INT PRIMARY KEY AUTO_INCREMENT,  		                -- 学生ID，主键
                              name VARCHAR(50) NOT NULL,  										-- 名
                              gender CHAR(1) NOT NULL,  											-- 性别
                              date_of_birth DATE NOT NULL,  										-- 出生日期
                              enrollment_date DATE NOT NULL,  									-- 入学日期
                              grade_offset INT DEFAULT 0,                                          -- 根据入学日期计算年级的偏移 默认值0
                              school_id INT NOT NULL,  											-- 学校ID
                              mobile VARCHAR(15) NOT NULL,                                         -- 手机号码
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,                      -- 创建时间
                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- 更新时间
                              deleted_at TIMESTAMP NULL DEFAULT NULL,                             -- 删除时间
                              FOREIGN KEY (school_id) REFERENCES usr_schools(school_id)               -- 外键，参考学校表
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = COMPACT;

-- 创建教师表
CREATE TABLE usr_teachers (
                              teacher_id INT PRIMARY KEY AUTO_INCREMENT, 	-- 教师ID，主键，自动递增
                              name VARCHAR(50) NOT NULL, 					-- 名，最多50个字符，不允许为空
                              gender CHAR(1) NOT NULL, 					-- 性别，1个字符，不允许为空
                              date_of_birth DATE NOT NULL, 				-- 出生日期，不允许为空
                              hire_date DATE NOT NULL, 					-- 雇佣日期，不允许为空
                              subject VARCHAR(50) NOT NULL, 				-- 教授科目，最多50个字符，不允许为空
                              school_id INT NOT NULL, 					-- 学校ID，不允许为空
                              mobile VARCHAR(15) NOT NULL,                 -- 手机号码
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- 更新时间
                              deleted_at TIMESTAMP NULL DEFAULT NULL,      -- 删除时间
                              FOREIGN KEY (school_id) REFERENCES usr_schools(school_id) -- 外键，参考学校表的school_id
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = COMPACT;

-- 创建班级表
CREATE TABLE usr_classes (
                             class_id INT PRIMARY KEY AUTO_INCREMENT,  			-- 班级ID，主键，自动递增
                             class_name VARCHAR(50) NOT NULL,  					-- 班级名称，最多50个字符，不允许为空
                             teacher_id INT NOT NULL,  							-- 教师ID，不允许为空
                             school_id INT NOT NULL,  							-- 学校ID，不允许为空
                             FOREIGN KEY (teacher_id) REFERENCES usr_teachers(teacher_id),  -- 外键，参考教师表
                             FOREIGN KEY (school_id) REFERENCES usr_schools(school_id)  -- 外键，参考学校表
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = COMPACT;

-- 创建学生班级关系表
CREATE TABLE usr_student_classes (
                                     student_class_id INT PRIMARY KEY AUTO_INCREMENT,  -- 学生班级关系ID，主键，自动递增
                                     student_id INT NOT NULL,  						-- 学生ID，不允许为空
                                     class_id INT NOT NULL,  						-- 班级ID，不允许为空
                                     enrollment_date DATE NOT NULL,  				-- 入班日期，不允许为空
                                     FOREIGN KEY (student_id) REFERENCES usr_students(student_id),  -- 外键，参考学生表
                                     FOREIGN KEY (class_id) REFERENCES usr_classes(class_id)  -- 外键，参考班级表
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = COMPACT;

-- 创建学生成绩表
CREATE TABLE usr_student_scores (
                                    score_id INT PRIMARY KEY AUTO_INCREMENT,  	-- 得分记录ID，主键，自动递增
                                    student_id INT NOT NULL,  					-- 学生ID，不允许为空
                                    year INT NOT NULL,  						-- 年，不允许为空
                                    semester INT NOT NULL,  					-- 学期，不允许为空（例如1或2）
                                    month INT NOT NULL,  						-- 月，不允许为空（1到12）
                                    moral_score DECIMAL(5, 2) NOT NULL DEFAULT 0.0,  	-- 德育成绩
                                    intellectual_score DECIMAL(5, 2) NOT NULL DEFAULT 0.0,  -- 智育成绩
                                    physical_score DECIMAL(5, 2) NOT NULL DEFAULT 0.0,  	-- 体育成绩
                                    aesthetic_score DECIMAL(5, 2) NOT NULL DEFAULT 0.0,  	-- 美育成绩
                                    labor_score DECIMAL(5, 2) NOT NULL DEFAULT 0.0,  		-- 劳育成绩
                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,       -- 创建时间
                                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- 更新时间
                                    deleted_at TIMESTAMP NULL DEFAULT NULL,               -- 删除时间
                                    FOREIGN KEY (student_id) REFERENCES usr_students(student_id) -- 外键，参考学生表的student_id
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = COMPACT;

-- 插入学校数据
INSERT INTO usr_schools (school_name, address) VALUES
                                                   ('Springfield High School', '742 Evergreen Terrace, Springfield'),
                                                   ('Shelbyville Elementary', '123 Main Street, Shelbyville');

-- 插入教师数据
INSERT INTO usr_teachers (name, gender, date_of_birth, hire_date, subject, school_id, mobile) VALUES
                                                                                                  ('John Smith', 'M', '1980-05-15', '2010-08-23', 'Mathematics', 1, '1234567890'),
                                                                                                  ('Jane Doe', 'F', '1985-07-20', '2012-09-15', 'Science', 2, '0987654321');

-- 插入班级数据
INSERT INTO usr_classes (class_name, teacher_id, school_id) VALUES
                                                                ('Math 101', 1, 1),
                                                                ('Science 101', 2, 2);

-- 插入学生数据
INSERT INTO usr_students (name, gender, date_of_birth, enrollment_date, grade_offset, school_id, mobile) VALUES
                                                                                                             ('Alice Johnson', 'F', '2005-03-12', '2020-09-01', 0, 1, '1112223333'),
                                                                                                             ('Bob Brown', 'M', '2006-06-24', '2020-09-01', 0, 2, '4445556666');

-- 插入学生班级关系数据
INSERT INTO usr_student_classes (student_id, class_id, enrollment_date) VALUES
                                                                            (1, 1, '2020-09-01'),
                                                                            (2, 2, '2020-09-01');

-- 插入学生成绩数据
INSERT INTO usr_student_scores (student_id, year, semester, month, moral_score, intellectual_score, physical_score, aesthetic_score, labor_score) VALUES
                                                                                                                                                      (1, 2020, 1, 9, 85.5, 90.0, 88.5, 92.0, 89.0),
                                                                                                                                                      (1, 2020, 2, 12, 87.0, 91.5, 89.5, 93.0, 90.0),
                                                                                                                                                      (2, 2020, 1, 9, 82.0, 85.5, 87.0, 88.0, 86.5),
                                                                                                                                                      (2, 2020, 2, 12, 84.0, 87.5, 88.0, 89.5, 88.0);
