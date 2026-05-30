import unittest

from predictor.lookup import (
    build_photo_key,
    build_sell_key,
    ensure_min_photos,
    normalize_car_part,
)


class LookupTest(unittest.TestCase):
    def test_normalize_car_part_replaces_spaces(self):
        self.assertEqual(normalize_car_part("Land Rover"), "Land_Rover")

    def test_build_photo_key_uses_normalized_make_and_model(self):
        self.assertEqual(
            build_photo_key("Land Rover", "Range Rover", 2020),
            "Land_Rover/Range_Rover/2020",
        )

    def test_build_sell_key_keeps_dataset_separator(self):
        self.assertEqual(build_sell_key("Toyota", "Camry"), "Toyota | Camry")

    def test_ensure_min_photos_repeats_first_photo(self):
        self.assertEqual(
            ensure_min_photos(["https://example.com/car.jpg"]),
            [
                "https://example.com/car.jpg",
                "https://example.com/car.jpg",
                "https://example.com/car.jpg",
            ],
        )

    def test_ensure_min_photos_keeps_empty_list_empty(self):
        self.assertEqual(ensure_min_photos([]), [])


if __name__ == "__main__":
    unittest.main()
