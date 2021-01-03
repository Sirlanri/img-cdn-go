DELETE FROM reply WHERE toUser NOT IN (
    SELECT userid FROM `user`
);

UPDATE picture SET picaddress=REPLACE(
    picaddress,"https://api.ri-co.cn/whisper/getpics","https://cdn.ri-co.cn/img"
);

UPDATE picture SET picaddress=REPLACE(
    picaddress,"https://api.ri-co.cn/whisper/getpics","https://cdn.ri-co.cn/img"
);

UPDATE user SET bannar=REPLACE(
    bannar,"https://api.ri-co.cn/whisper/getpics","https://cdn.ri-co.cn/img"
);