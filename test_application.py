import os
import subprocess
import time
import unittest
from typing import Optional

import requests


class TestCalculatorServices(unittest.TestCase):
    """Интеграционный тест."""
    auth_service_url = "http://localhost:8081"
    main_service_url = "http://localhost:8080"
    test_username = "testuser"
    test_password = "testpass123"

    def test_full_flow(self):
        """Полный тест: регистрация -> авторизация -> вычисление выражения"""
        # 1. Регистрация нового пользователя
        register_response = requests.post(
            f"{self.auth_service_url}/register",
            json={"username": self.test_username, "password": self.test_password},
            timeout=5
        )
        if register_response.status_code != 400:
            self.assertEqual(register_response.status_code, 201, "Регистрация не удалась")

        # 2. Авторизация (получаем JWT токен)
        login_response = requests.post(
            f"{self.auth_service_url}/login",
            json={"username": self.test_username, "password": self.test_password},
            timeout=5
        )
        self.assertEqual(login_response.status_code, 200, "Авторизация не удалась")
        token = login_response.json().get("token")
        self.assertIsNotNone(token, "Токен не получен")

        headers = {"Authorization": f"{token}"}

        # 3. Отправляем выражение на вычисление
        expression = "2+2*2"
        calc_response = requests.post(
            f"{self.main_service_url}/api/v1/calculate",
            json={"expression": expression},
            headers=headers,
            timeout=5
        )
        self.assertEqual(calc_response.status_code, 201, "Ошибка при отправке выражения")
        expression_id = calc_response.json().get("id")
        self.assertIsNotNone(expression_id, "ID выражения не получен")

        # 4. Проверяем статус вычисления (может потребоваться несколько попыток)
        max_attempts = 10
        for _ in range(max_attempts):
            status_response = requests.get(
                f"{self.main_service_url}/api/v1/expression/{expression_id}",
                headers=headers,
                timeout=5
            )
            self.assertEqual(status_response.status_code, 200, "Ошибка при проверке статуса")
            
            status = status_response.json().get("expression", {}).get("status")
            result = status_response.json().get("expression", {}).get("result")

            if status == "solved":
                self.assertEqual(float(result), 6, "Неправильный результат вычисления")
                break
            
            time.sleep(1)  # Ждем секунду перед следующей проверкой
        else:
            self.fail("Вычисление не завершилось за отведенное время")

        # 5. Проверяем список всех выражений
        list_response = requests.get(
            f"{self.main_service_url}/api/v1/expressions",
            headers=headers,
            timeout=5
        )
        self.assertEqual(list_response.status_code, 200)
        expressions = list_response.json().get("expressions", [])
        self.assertGreater(len(expressions), 0, "Список выражений пуст")


if __name__ == "__main__":
    unittest.main()