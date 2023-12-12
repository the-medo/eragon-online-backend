CREATE VIEW view_worlds AS
SELECT
    w.*,
    vm.*
FROM
    worlds w
    JOIN view_modules vm ON w.id = vm.module_world_id
;