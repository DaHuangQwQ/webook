CREATE DATABASE dahuang;

INSERT INTO `sys_auth_rules` (id, p_id, name, title, icon, link_url, redirect, menu_type, weigh, is_iframe, path, component, is_cached, remark, is_affix, is_link, module_type, `condition`)
VALUES
    (1, 0, 'api/v1/system/auth', '权限管理', 'ele-Stamp', '', '', 0, 30, 0, '/system/auth', 'layout/routerView/parent', 0, '', 0, 0, '0', ''),
    (2, 1, 'api/v1/system/auth/menuList', '菜单管理', 'ele-Calendar', '', '', 1, 0, 0, '/system/auth/menuList', 'system/menu/index', 0, '', 0, 0, '', ''),
    (3, 2, 'api/v1/system/menu/add', '添加菜单', '', '', '', 2, 0, 0, '', '', 0, '', 0, 0, '', ''),
    (4, 2, 'api/v1/system/menu/update', '修改菜单', '', '', '', 2, 0, 0, '', '', 0, '', 0, 0, '', ''),
    (10, 1, 'api/v1/system/role/list', '角色管理', 'iconfont icon-juxingkaobei', '', '', 1, 0, 0, '/system/auth/roleList', 'system/role/index', 0, '', 0, 0, '', ''),
    (11, 2, 'api/v1/system/menu/delete', '删除菜单', '', '', '', 2, 0, 0, '', '', 0, '', 0, 0, '', ''),
    (12, 10, 'api/v1/system/role/add', '添加角色', '', '', '', 2, 0, 0, '', '', 0, '', 0, 0, '', ''),
    (13, 10, '/api/v1/system/role/edit', '修改角色', '', '', '', 2, 0, 0, '', '', 0, '', 0, 0, '', ''),
    (14, 10, '/api/v1/system/role/delete', '删除角色', '', '', '', 2, 0, 0, '', '', 0, '', 0, 0, '', ''),
    (15, 1, 'api/v1/system/dept/list', '部门管理', 'iconfont icon-siweidaotu', '', '', 1, 0, 0, '/system/auth/deptList', 'system/dept/index', 0, '', 0, 0, '', ''),
    (16, 17, 'aliyun', '阿里云-iframe', 'iconfont icon-diannao1', '', '', 1, 0, 0, '/demo/outLink/aliyun', 'layout/routerView/iframes', 1, '', 0, 1, '', 'https://www.aliyun.com/daily-act/ecs/activity_selection?spm=5176.8789780.J_3965641470.5.568845b58KHj51'),
    (17, 0, 'outLink', '外链测试', 'iconfont icon-zhongduancanshu', '', '', 0, 20, 0, '/demo/outLink', 'layout/routerView/parent', 0, '', 0, 0, '', ''),
    (18, 17, 'tenyun', '腾讯云-外链', 'iconfont icon-shouye_dongtaihui', '', '', 1, 0, 0, '/demo/outLink/tenyun', 'layout/routerView/link', 1, '', 0, 0,  '', 'https://cloud.tencent.com/act/new?cps_key=20b1c3842f74986b2894e2c5fcde7ea2&fromSource=gwzcw.3775555.3775555.3775555&utm_id=gwzcw.3775555.3775555.3775555&utm_medium=cpc'),
    (19, 15, 'api/v1/system/dept/add', '添加部门', '', '', '', 2, 0, 0, '', '', 0, '', 0, 0, '', ''),
    (20, 15, 'api/v1/system/dept/edit', '修改部门', '', '', '', 2, 0, 0, '', '', 0, '', 0, 0, '', ''),
    (21, 15, 'api/v1/system/dept/delete', '删除部门', '', '', '', 2, 0, 0, '', '', 0, '', 0, 0, '', ''),
    (22, 1, 'api/v1/system/post/list', '岗位管理', 'iconfont icon-neiqianshujuchucun', '', '', 1, 0, 0, '/system/auth/postList', 'system/post/index', 0, '', 0, 0, '', ''),
    (23, 22, 'api/v1/system/post/add', '添加岗位', '', '', '', 2, 0, 0, '', '', 0, '', 0, 0, '', ''),
    (24, 22, 'api/v1/system/post/edit', '修改岗位', '', '', '', 2, 0, 0, '', '', 0, '', 0, 0, '', ''),
    (25, 22, 'api/v1/system/post/delete', '删除岗位', '', '', '', 2, 0, 0, '', '', 0, '', 0, 0, '', ''),
    (26, 1, 'api/v1/system/user/list', '用户管理', 'ele-User', '', '', 1, 0, 0, '/system/auth/user/list', 'system/user/index', 0, '', 0, 0, '', ''),
    (27, 0, 'api/v1/system/dict', '系统配置', 'iconfont icon-shuxingtu', '', '', 0, 40, 0, '/system/dict', 'layout/routerView/parent', 0, '', 0, 0, '654', ''),
    (28, 27, 'api/v1/system/dict/type/list', '字典管理', 'iconfont icon-crew_feature', '', '', 1, 0, 0, '/system/dict/type/list', 'system/dict/index', 0, '', 0, 0, '', ''),
    (29, 27, 'api/v1/system/dict/dataList', '字典数据管理', 'iconfont icon-putong', '', '', 1, 0, 1, '/system/dict/data/list/:dictType', 'system/dict/dataList', 0, '', 0, 0, '', ''),
    (30, 27, 'api/v1/system/config/list', '参数管理', 'ele-Cherry', '', '', 1, 0, 0, '/system/config/list', 'system/config/index', 0, '', 0, 0, '', ''),
    (31, 0, 'api/v1/system/monitor', '系统监控', 'iconfont icon-xuanzeqi', '', '', 0, 30, 0, '/system/monitor', 'layout/routerView/parent', 0, '', 0, 0, '', ''),
    (32, 31, 'api/v1/system/monitor/server', '服务监控', 'iconfont icon-shuju', '', '', 1, 0, 0, '/system/monitor/server', 'system/monitor/server/index', 0, '', 0, 0, '', ''),
    (33, 35, 'api/swagger', 'api文档', 'iconfont icon--chaifenlie', '', '', 1, 0, 0, '/system/swagger', 'layout/routerView/iframes', 1, '', 0, 1, '', 'http://localhost:8808/swagger'),
    (34, 31, 'api/v1/system/loginLog/list', '登录日志', 'ele-Finished', '', '', 1, 0, 0, '/system/monitor/loginLog', 'system/monitor/loginLog/index', 0, '', 0, 0, '', ''),
    (35, 0, 'api/v1/system/tools', '系统工具', 'iconfont icon-zujian', '', '', 0, 25, 0, '/system/tools', 'layout/routerView/parent', 0, '', 0, 0, '', ''),
    (38, 31, 'api/v1/system/operLog/list', '操作日志', 'iconfont icon-bolangnengshiyanchang', '', '', 1, 0, 0, '/system/monitor/operLog', 'system/monitor/operLog/index', 0, '', 0, 0, '', '')



-- -- 创建学校表
-- CREATE TABLE usr_schools (
--                              school_id INT PRIMARY KEY AUTO_INCREMENT, 	-- 学校ID，主键，自动递增
--                              school_name VARCHAR(100) NOT NULL, 					-- 学校名称，最多100个字符，不允许为空
--                              address VARCHAR(255) NOT NULL 							-- 学校地址，最多255个字符，不允许为空
-- ) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = COMPACT;
--
-- -- 创建学生表
-- CREATE TABLE usr_students (
--                               student_id INT PRIMARY KEY AUTO_INCREMENT,  		                -- 学生ID，主键
--                               name VARCHAR(50) NOT NULL,  										-- 名
--                               gender CHAR(1) NOT NULL,  											-- 性别
--                               date_of_birth DATE NOT NULL,  										-- 出生日期
--                               enrollment_date DATE NOT NULL,  									-- 入学日期
--                               grade_offset INT DEFAULT 0,                                          -- 根据入学日期计算年级的偏移 默认值0
--                               school_id INT NOT NULL,  											-- 学校ID
--                               mobile VARCHAR(15) NOT NULL,                                         -- 手机号码
--                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,                      -- 创建时间
--                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- 更新时间
--                               deleted_at TIMESTAMP NULL DEFAULT NULL,                             -- 删除时间
--                               FOREIGN KEY (school_id) REFERENCES usr_schools(school_id)               -- 外键，参考学校表
-- ) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = COMPACT;
--
-- -- 创建教师表
-- CREATE TABLE usr_teachers (
--                               teacher_id INT PRIMARY KEY AUTO_INCREMENT, 	-- 教师ID，主键，自动递增
--                               name VARCHAR(50) NOT NULL, 					-- 名，最多50个字符，不允许为空
--                               gender CHAR(1) NOT NULL, 					-- 性别，1个字符，不允许为空
--                               date_of_birth DATE NOT NULL, 				-- 出生日期，不允许为空
--                               hire_date DATE NOT NULL, 					-- 雇佣日期，不允许为空
--                               subject VARCHAR(50) NOT NULL, 				-- 教授科目，最多50个字符，不允许为空
--                               school_id INT NOT NULL, 					-- 学校ID，不允许为空
--                               mobile VARCHAR(15) NOT NULL,                 -- 手机号码
--                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
--                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- 更新时间
--                               deleted_at TIMESTAMP NULL DEFAULT NULL,      -- 删除时间
--                               FOREIGN KEY (school_id) REFERENCES usr_schools(school_id) -- 外键，参考学校表的school_id
-- ) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = COMPACT;
--
-- -- 创建班级表
-- CREATE TABLE usr_classes (
--                              class_id INT PRIMARY KEY AUTO_INCREMENT,  			-- 班级ID，主键，自动递增
--                              class_name VARCHAR(50) NOT NULL,  					-- 班级名称，最多50个字符，不允许为空
--                              teacher_id INT NOT NULL,  							-- 教师ID，不允许为空
--                              school_id INT NOT NULL,  							-- 学校ID，不允许为空
--                              FOREIGN KEY (teacher_id) REFERENCES usr_teachers(teacher_id),  -- 外键，参考教师表
--                              FOREIGN KEY (school_id) REFERENCES usr_schools(school_id)  -- 外键，参考学校表
-- ) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = COMPACT;
--
-- -- 创建学生班级关系表
-- CREATE TABLE usr_student_classes (
--                                      student_class_id INT PRIMARY KEY AUTO_INCREMENT,  -- 学生班级关系ID，主键，自动递增
--                                      student_id INT NOT NULL,  						-- 学生ID，不允许为空
--                                      class_id INT NOT NULL,  						-- 班级ID，不允许为空
--                                      enrollment_date DATE NOT NULL,  				-- 入班日期，不允许为空
--                                      FOREIGN KEY (student_id) REFERENCES usr_students(student_id),  -- 外键，参考学生表
--                                      FOREIGN KEY (class_id) REFERENCES usr_classes(class_id)  -- 外键，参考班级表
-- ) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = COMPACT;
--
-- -- 创建学生成绩表
-- CREATE TABLE usr_student_scores (
--                                     score_id INT PRIMARY KEY AUTO_INCREMENT,  	-- 得分记录ID，主键，自动递增
--                                     student_id INT NOT NULL,  					-- 学生ID，不允许为空
--                                     year INT NOT NULL,  						-- 年，不允许为空
--                                     semester INT NOT NULL,  					-- 学期，不允许为空（例如1或2）
--                                     month INT NOT NULL,  						-- 月，不允许为空（1到12）
--                                     moral_score DECIMAL(5, 2) NOT NULL DEFAULT 0.0,  	-- 德育成绩
--                                     intellectual_score DECIMAL(5, 2) NOT NULL DEFAULT 0.0,  -- 智育成绩
--                                     physical_score DECIMAL(5, 2) NOT NULL DEFAULT 0.0,  	-- 体育成绩
--                                     aesthetic_score DECIMAL(5, 2) NOT NULL DEFAULT 0.0,  	-- 美育成绩
--                                     labor_score DECIMAL(5, 2) NOT NULL DEFAULT 0.0,  		-- 劳育成绩
--                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,       -- 创建时间
--                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- 更新时间
--                                     deleted_at TIMESTAMP NULL DEFAULT NULL,               -- 删除时间
--                                     FOREIGN KEY (student_id) REFERENCES usr_students(student_id) -- 外键，参考学生表的student_id
-- ) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = COMPACT;
--
-- -- 插入学校数据
-- INSERT INTO usr_schools (school_name, address) VALUES
--                                                    ('Springfield High School', '742 Evergreen Terrace, Springfield'),
--                                                    ('Shelbyville Elementary', '123 Main Street, Shelbyville');
--
-- -- 插入教师数据
-- INSERT INTO usr_teachers (name, gender, date_of_birth, hire_date, subject, school_id, mobile) VALUES
--                                                                                                   ('John Smith', 'M', '1980-05-15', '2010-08-23', 'Mathematics', 1, '1234567890'),
--                                                                                                   ('Jane Doe', 'F', '1985-07-20', '2012-09-15', 'Science', 2, '0987654321');
--
-- -- 插入班级数据
-- INSERT INTO usr_classes (class_name, teacher_id, school_id) VALUES
--                                                                 ('Math 101', 1, 1),
--                                                                 ('Science 101', 2, 2);
--
-- -- 插入学生数据
-- INSERT INTO usr_students (name, gender, date_of_birth, enrollment_date, grade_offset, school_id, mobile) VALUES
--                                                                                                              ('Alice Johnson', 'F', '2005-03-12', '2020-09-01', 0, 1, '1112223333'),
--                                                                                                              ('Bob Brown', 'M', '2006-06-24', '2020-09-01', 0, 2, '4445556666');
--
-- -- 插入学生班级关系数据
-- INSERT INTO usr_student_classes (student_id, class_id, enrollment_date) VALUES
--                                                                             (1, 1, '2020-09-01'),
--                                                                             (2, 2, '2020-09-01');
--
-- -- 插入学生成绩数据
-- INSERT INTO usr_student_scores (student_id, year, semester, month, moral_score, intellectual_score, physical_score, aesthetic_score, labor_score) VALUES
--                                                                                                                                                       (1, 2020, 1, 9, 85.5, 90.0, 88.5, 92.0, 89.0),
--                                                                                                                                                       (1, 2020, 2, 12, 87.0, 91.5, 89.5, 93.0, 90.0),
--                                                                                                                                                       (2, 2020, 1, 9, 82.0, 85.5, 87.0, 88.0, 86.5),
--                                                                                                                                                       (2, 2020, 2, 12, 84.0, 87.5, 88.0, 89.5, 88.0);
