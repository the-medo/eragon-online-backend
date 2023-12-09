CREATE VIEW view_worlds AS
SELECT
    w.*,
    vm.id as module_id,
    vm.menu_id,
    vm.header_img_id,
    vm.thumbnail_img_id,
    vm.avatar_img_id,
    vm.tags
FROM
    worlds w
        JOIN view_modules vm ON w.id = vm.world_id
;