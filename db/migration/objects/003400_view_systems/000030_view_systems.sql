CREATE VIEW view_systems AS
SELECT
    s.*,
    vm.id as module_id,
    vm.menu_id,
    vm.header_img_id,
    vm.thumbnail_img_id,
    vm.avatar_img_id,
    vm.tags
FROM
    systems s
    JOIN view_modules vm ON s.id = vm.system_id
;