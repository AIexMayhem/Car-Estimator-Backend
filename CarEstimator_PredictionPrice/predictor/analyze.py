import os
import joblib
import pandas as pd
import numpy as np
import random
from scipy.interpolate import make_interp_spline
import matplotlib.pyplot as plt
from io import BytesIO

plt.switch_backend('agg')

BASE_DIR = os.path.dirname(os.path.abspath(__file__))
DATA_DIR = os.path.join(BASE_DIR, 'model_data')
MODEL_PATH = os.path.join(DATA_DIR, 'linear_model.sav')
DUMMIES = os.path.join(DATA_DIR, 'car_dummies.csv')
PICS_PATH = os.path.join(DATA_DIR, 'pictures.csv')
SELLS_PATH = os.path.join(DATA_DIR, 'sells.csv')

load_model = joblib.load(open(MODEL_PATH, 'rb'))
FEATURE_NAMES = list(load_model.feature_names_in_)

pics = pd.read_csv(PICS_PATH)
sells = pd.read_csv(SELLS_PATH)

def get_car_info(make: str, model: str, year: int, hp: int,
                 body: str, yearsell: int, odometer: int, color: str) -> int:
    dummy = dict.fromkeys(FEATURE_NAMES, 0)

    for k, v in {
        'Year': year,
        'HP': hp,
        'Odometer': odometer,
        'Yearsell': yearsell
    }.items():
        if k in dummy:
            dummy[k] = v

    for prefix, val in (
        ('Make_', make.replace(' ', '_')),
        ('Model_', model.replace(' ', '_')),
        ('Body_', body),
        ('Color_', color)
    ):
        col = f"{prefix}{val}"
        if col in dummy:
            dummy[col] = 1

    df = pd.DataFrame([dummy], columns=FEATURE_NAMES)
    y_pred = load_model.predict(df)[0]
    price  = max(int(round(y_pred, 0)), 0)

    return price + random.randint(-price // 10, price // 10)

def get_photos(make: str, model: str, year: int) -> list[str]:
    key = f"{make.replace(' ', '_')}/{model.replace(' ', '_')}/{year}"
    for _, row in pics.iterrows():
        if row['Car'] == key:
            photos = row['Pics'].split()
            while 0 < len(photos) < 3:
                photos.append(photos[0])
            return photos
    return []

def get_sells(make: str, model: str) -> int:
    key = f"{make} | {model}"
    for _, row in sells.iterrows():
        if row['Car'] == key:
            return int(row['Count'])
    return 0

def graph_build(make: str, model: str, year: int, hp: int,
                body: str, yearsell: int, odometer: int, color: str) -> bytes:

    period = 8
    xs, ys = [], []
    base_year = yearsell - 4

    for i in range(period):
        y = base_year + i + 1
        price = get_car_info(make, model, year, hp, body, y, odometer, color)
        xs.append(y)
        ys.append(price if price > 0 else 0)

    xnew = np.linspace(min(xs), max(xs), 300)
    spl = make_interp_spline(xs, ys, k=3)
    smooth = spl(xnew)

    plt.clf()
    plt.plot(xnew, smooth)
    plt.box(False)

    buf = BytesIO()
    plt.savefig(buf, format='png')
    buf.seek(0)
    return buf.read()