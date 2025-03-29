# isolate-wrapper

**isolate-wrapper** is the core evaluation engine for the Patito Virtual Judge.  
It listens to a RabbitMQ queue, retrieves submitted code, runs it safely using `isolate`, and sends the result back via an async callback.

---

## JSON Task Format


```json
{
  "id": "123",
  "uniq_id": "a1b2c3d4e5",
  "time_submit_code": "2025-03-28T14:23:00Z",
  "problem_id": 101,
  "site_id": 2,
  "code": "#include <iostream>\nint main() { std::cout << \"Hola Patito!\"; return 0; }",
  "language": "cpp",
  "callback": "https://patito.site/api/v1/submit/callback/123",
  "run_limits": {
    "time": 2,
    "memory": 65536,
    "output": 1024
  }
}

![Architecture](image.png "Architecture")
