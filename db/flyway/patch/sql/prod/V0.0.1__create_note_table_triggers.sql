CREATE TRIGGER trigger_check_one_actual_note
BEFORE INSERT
   ON note FOR EACH ROW
BEGIN
    DECLARE actualRowCount INT;

    SELECT COUNT(*) 
    INTO actualRowCount
    FROM note 
    WHERE note_guid = new.note_guid and actual = 1;

    if (new.actual = 1) then
        if actualRowCount >= 1 then
            -- signal sqlstate '45000' set message_text = concat('Нельзя добавить две актуальные записи для note', new.note_guid);
            signal sqlstate '45000' set message_text = 'Нельзя добавить две актуальные записи для note';
        end if;
    end if;
END;