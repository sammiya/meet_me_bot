import unittest
from datetime import datetime
from util import parse_next_date

class TestUtil(unittest.TestCase):
    def test_parse_next_date(self):
        # テストケースを定義
        test_cases = [
            ("次回は 3/7 ね。", datetime(datetime.now().year, 3, 7, 21)),
            ("次回12/5", datetime(datetime.now().year, 12, 5, 21)),
            ("次回 5/13（金）", datetime(datetime.now().year, 5, 13, 21)),
            ("foo", None),  # 日付が含まれない場合はNoneを返す
        ]

        # テストを実行
        for text, expected in test_cases:
            with self.subTest(text=text):
                self.assertEqual(parse_next_date(text), expected)

if __name__ == '__main__':
    unittest.main()
