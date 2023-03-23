import discord
import os
from dotenv import load_dotenv
import datetime
from discord.ext import commands, tasks
import pytz
import asyncio
import sqlite3
import locale

from util import parse_next_date, is_asking_next_date

# 日本語曜日名を定義
WEEKDAYS = ['月', '火', '水', '木', '金', '土', '日']

# .env ファイルから環境変数を読み込む
load_dotenv()

# Discord Bot トークンを取得する
TOKEN = os.getenv('DISCORD_TOKEN')
CHANNEL_ID = int(os.getenv('DISCORD_CHANNEL'))
MEETING_ANNOUNCER_ID = int(os.getenv('DISCORD_MEETING_ANNOUNCER'))

# タイムゾーンを設定
JST = pytz.timezone('Asia/Tokyo')

intents = discord.Intents(messages=True, guilds=True, reactions=True, message_content=True)

bot = commands.Bot(command_prefix='!', intents=intents)

# SQLiteに接続
conn = sqlite3.connect('meet-me-bot.db')
cursor = conn.cursor()

@tasks.loop(seconds=60)  # 60秒ごとに実行
async def scheduled_message():
    cursor.execute("SELECT date, pre_notification_sent FROM meetings")
    row = cursor.fetchone()

    if row is None:
        return

    scheduled_time_str, pre_notification_sent = row
    scheduled_time = datetime.datetime.strptime(row[0], '%Y-%m-%d %H:%M:%S').astimezone(JST)
    now_time = datetime.datetime.now(JST)
    delta_time = scheduled_time - now_time
    remaining_seconds = delta_time.total_seconds()

    if not pre_notification_sent and 3300 < remaining_seconds < 3660:
        remaining_for_prenotification = remaining_seconds - 3600
        await asyncio.sleep(remaining_for_prenotification)  # 差分ぶんだけ待機

        channel = bot.get_channel(CHANNEL_ID)
        await channel.send('1時間前だよ')
        cursor.execute("UPDATE meetings SET pre_notification_sent = 1")
        conn.commit()

    if remaining_seconds < 60:
        await asyncio.sleep(remaining_seconds)  # 差分ぶんだけ待機

        channel = bot.get_channel(CHANNEL_ID)
        await channel.send('はじまるよ')
        cursor.execute("DELETE FROM meetings")
        conn.commit()

# メッセージが送信されたときのイベントを処理する関数を定義する
@bot.event
async def on_message(message):
    if is_asking_next_date(message.content):
        channel = bot.get_channel(CHANNEL_ID)
        cursor.execute("SELECT date FROM meetings")
        row = cursor.fetchone()
        if row is None:
            await channel.send('次回予定は未定だよ')
            return

        scheduled_time_str = row[0]
        scheduled_time = datetime.datetime.strptime(scheduled_time_str, '%Y-%m-%d %H:%M:%S').astimezone(JST)
        formatted = scheduled_time.strftime('%-m/%-d') + f'（{WEEKDAYS[scheduled_time.weekday()]}）'
        await channel.send(f'次回予定: {formatted}')

    if (message.author.id != MEETING_ANNOUNCER_ID):
        return

    parsed = parse_next_date(message.content)
    if (parsed is None):
        return

    try:
        conn.execute('BEGIN')
        conn.execute('DELETE FROM meetings')
        cursor.execute(f"INSERT INTO meetings(date, pre_notification_sent) VALUES ('{parsed.strftime('%Y-%m-%d %H:%M:%S')}', 0)")
        conn.commit()
        await message.add_reaction("👍")
    except sqlite3.IntegrityError as e:
        print(f"Error inserting parsed date: {e}")
        await message.add_reaction("👎")
        conn.rollback()

@bot.event
async def on_ready():
    """
    BotがDiscordにログインして準備ができたときに呼び出される関数。
    ログインしたことを確認するために、Botの名前を出力する。
    """
    print(f'We have logged in as {bot.user}')
    channel = bot.get_channel(CHANNEL_ID)
    # await channel.send('Hello!')
    scheduled_message.start()

bot.run(TOKEN)
