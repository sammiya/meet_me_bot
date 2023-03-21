import re
from datetime import datetime, timedelta

def parse_next_date(text):
    match = re.search(r'次回[^\d]*(\d{1,2}/\d{1,2})', text)
    if not match:
        return None

    # 日付を抽出する
    date_str = str(datetime.now().year) + "/" + match.group(1) + " 21:00"

    # 日付オブジェクトを作成する
    try:
        next_date = datetime.strptime(date_str, '%Y/%m/%d %H:%M')
    except ValueError:
        return None

    return next_date

def is_asking_next_date(text):
    pattern = r'次.*いつ'
    return re.search(pattern, text) is not None
