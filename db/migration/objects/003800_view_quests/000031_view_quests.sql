CREATE VIEW view_quests AS
SELECT
    q.*,
    qs.status as status,
    qs.can_join as can_join,
    vm.id as module_id,
    vm.menu_id,
    vm.header_img_id,
    vm.thumbnail_img_id,
    vm.avatar_img_id,
    vm.tags
FROM
    quests q
    JOIN view_modules vm ON q.id = vm.quest_id
    JOIN quest_settings qs ON q.id = qs.quest_id
;