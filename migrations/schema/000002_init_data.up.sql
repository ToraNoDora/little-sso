insert into users (
    username,
    email,
    pass_hash
) values
    (
        'default_username_1',
        'default_test_mail_1@mail.ru',
        E'\\x' || sha256('qwerty123')
    ),
    (
        'default_username_2',
        'default_test_mail_2@mail.ru',
        E'\\x' || sha256('qwerty123')),
    (
        'default_username_3',
        'default_test_mail_3@mail.ru',
        E'\\x' || sha256('qwerty123')
    );

insert into apps (
    name,
    description,
    secret
) values
    ('default_test_app_1', 'default test description of app 1', 'test-secret-test-app-1'),
    ('default_test_app_2', 'default test description of app 2', 'test-secret-test-app-2') on conflict do nothing;

insert into admins (
    user_id,
    app_id
) select u.id, a.id from users as u join apps as a on a.name = 'default_test_app_1' where u.email = 'default_test_mail_1@mail.ru' on conflict do nothing;

insert into admins (
    user_id,
    app_id,
    is_admin
) select u.id, a.id, true from users as u join apps as a on a.name = 'default_test_app_2' where u.email = 'default_test_mail_2@mail.ru' on conflict do nothing;

insert into groups (
    app_id,
    name,
    description
) select a.id, 'default_test_group_1', 'default test description 1' from apps as a where a.name = 'default_test_app_2' on conflict do nothing;

insert into groups (
    app_id,
    name,
    description
) select a.id, 'default_test_group_2', 'default test description 2' from apps as a where a.name = 'default_test_app_2' on conflict do nothing;

insert into roles (
    name,
    description
) values
    ('default_test_role_1', 'default test role description 1'),
    ('default_test_role_2', 'default test role description 2') on conflict do nothing;

insert into groups_roles (
    group_id,
    role_id
) select
    g.id, r.id
from groups as g join roles as r on r.name = 'default_test_role_1'
where g.name = 'default_test_group_1' on conflict do nothing;

insert into groups_roles (
    group_id,
    role_id
) select
    g.id, r.id
from groups as g join roles as r on r.name = 'default_test_role_2'
where g.name = 'default_test_group_2' on conflict do nothing;

insert into users_permissions (
    user_id,
    group_id
) select u.id, g.id from users as u join groups as g on g.name = 'default_test_group_1' where u.email = 'default_test_mail_2@mail.ru' on conflict do nothing;

insert into users_permissions (
    user_id,
    group_id,
    add_flag
) select u.id, g.id, true from users as u join groups as g on g.name = 'default_test_group_2' where u.email = 'default_test_mail_2@mail.ru' on conflict do nothing;

