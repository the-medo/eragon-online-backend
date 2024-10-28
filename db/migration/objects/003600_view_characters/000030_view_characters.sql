CREATE VIEW view_characters AS
SELECT
    c.*,
    vm.id as module_id,
    vm.menu_id,
    vm.header_img_id,
    vm.thumbnail_img_id,
    vm.avatar_img_id,
    vm.tags
FROM
    characters c
    JOIN view_modules vm ON c.id = vm.character_id
;