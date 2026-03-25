
INSERT INTO users (id, name, email, password, status, reserved)
VALUES ('1000', 'Administrator', 'admin@ekhoes.local', ?, 'enabled', 1);

insert into user_roles("user_id", "roles") values ('1000', 'ADMIN');
