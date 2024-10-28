CREATE VIEW view_quests AS
SELECT
    q.*,
    vm.id as module_id,
    vm.menu_id,
    vm.header_img_id,
    vm.thumbnail_img_id,
    vm.avatar_img_id,
    vm.tags
FROM
    quests q
    JOIN view_modules vm ON q.id = vm.quest_id
;