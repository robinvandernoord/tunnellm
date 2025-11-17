
import requests
import json
response = requests.get(
  url="https://openrouter.ai/api/v1/key",
  headers={
    "Authorization": f"Bearer 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"
  }
)
print(json.dumps(response.json(), indent=2))

