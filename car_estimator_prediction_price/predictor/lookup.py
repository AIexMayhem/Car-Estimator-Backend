def normalize_car_part(value: str) -> str:
    return value.replace(" ", "_")


def build_photo_key(make: str, model: str, year: int) -> str:
    return f"{normalize_car_part(make)}/{normalize_car_part(model)}/{year}"


def build_sell_key(make: str, model: str) -> str:
    return f"{make} | {model}"


def ensure_min_photos(photos: list[str], min_count: int = 3) -> list[str]:
    if not photos:
        return []

    result = list(photos)
    while len(result) < min_count:
        result.append(result[0])
    return result
