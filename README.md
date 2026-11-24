# GDG Backend Side Project — 식당 웨이팅 서비스

Go로 구현한 **식당 웨이팅(대기 등록) 관리 서비스**입니다.
손님은 로그인 없이 식당에 웨이팅을 등록하고, 사장님은 로그인해서 자기 식당의 대기 목록을 확인하고 입장/취소 처리를 합니다.

---

## 목차
1. [주요 기능](#주요-기능)
2. [기술 스택](#기술-스택)
3. [빠른 시작](#빠른-시작)
4. [환경 변수](#환경-변수)
5. [아키텍처](#아키텍처)
6. [데이터베이스](#데이터베이스)
7. [인증과 권한](#인증과-권한)
8. [API 명세](#api-명세)
9. [화면](#화면)
10. [테스트](#테스트)
11. [알려진 한계와 다음 단계](#알려진-한계와-다음-단계)
12. [기여자](#기여자)

---

## 주요 기능

**손님**
- 식당 목록과 상세 정보 조회 (현재 대기 팀 수, 예상 대기 시간)
- 인원수와 휴대폰 번호로 웨이팅 등록 (로그인 불필요)

**사장님**
- 회원가입 / 로그인 (JWT 발급)
- 내 계정을 식당의 사장님으로 등록
- 내 식당의 웨이팅 목록 조회, 입장·취소 처리
- 내 정보 조회 / 수정 / 탈퇴

**공통**
- 통일된 JSON 에러 응답 `{"error": "..."}`
- 입력 검증 (이메일, 비밀번호 길이, 휴대폰 번호 형식, 인원수 범위, 중복 웨이팅 방지)
- 요청 로깅, panic 복구, CORS, IP별 요청 속도 제한
- graceful shutdown

---

## 기술 스택

| 구분 | 사용 기술 |
|---|---|
| 언어 | Go 1.23 |
| 라우터 | gorilla/mux |
| DB | MySQL 8 (`go-sql-driver/mysql`, ORM 없이 SQL 직접 작성) |
| 인증 | JWT HS256 (`golang-jwt/jwt/v5`), bcrypt (`golang.org/x/crypto`) |
| 요청 제한 | `golang.org/x/time/rate` (토큰 버킷) |
| 설정 | 환경 변수 + `.env` (`joho/godotenv`) |
| 인프라 | Docker, Docker Compose |
| 개발 도구 | air (코드 저장 시 자동 재빌드) |

---

## 빠른 시작

### 방법 1. Docker Compose (권장)

```bash
cp .env.example .env
```

`.env`에서 `JWT_SECRET`을 길고 무작위인 값으로 바꾼 뒤 실행합니다.

```bash
docker compose up --build
```

- MySQL 컨테이너가 처음 뜰 때 `db/schema.sql`(테이블)과 `db/seed.sql`(샘플 식당 3개)이 자동으로 적용됩니다.
- 앱은 MySQL이 준비된 뒤(healthcheck) 시작되며, http://localhost:8080 에서 접속할 수 있습니다.
- 스키마를 처음부터 다시 적용하려면 볼륨을 지우고 다시 띄웁니다.

```bash
docker compose down -v
```

### 방법 2. 로컬 실행

MySQL 8이 설치되어 있어야 합니다.

```bash
cp .env.example .env
```

```bash
mysql -u root -p -e "CREATE DATABASE gdgdb CHARACTER SET utf8mb4"
```

```bash
mysql -u root -p gdgdb < db/schema.sql
```

```bash
mysql -u root -p gdgdb < db/seed.sql
```

```bash
go run .
```

개발 중에는 `air`를 쓰면 저장할 때마다 자동으로 다시 빌드됩니다 (`air.toml` 참고).

### 동작 확인

```bash
curl http://localhost:8080/api/healthz
```

```json
{"status":"ok"}
```

---

## 환경 변수

`.env.example`에 전체 예시가 있습니다.

| 변수 | 기본값 | 설명 |
|---|---|---|
| `PORT` | `8080` | 서버 포트 |
| `DB_HOST` | `localhost` | MySQL 호스트 (Compose에서는 `mysql`로 자동 설정) |
| `DB_PORT` | `3306` | MySQL 포트 |
| `DB_USER` | `root` | MySQL 사용자 |
| `DB_PASSWORD` | (빈 값) | MySQL 비밀번호 |
| `DB_NAME` | `gdgdb` | 데이터베이스 이름 |
| `JWT_SECRET` | **필수** | JWT 서명 키. 16자 미만이면 서버가 시작되지 않음 |
| `JWT_TTL` | `24h` | 토큰 유효 기간 |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:8080` | 허용 Origin, 쉼표로 구분. `*`는 전체 허용 |
| `RATE_LIMIT_RPS` | `10` | IP당 초당 허용 요청 수 |
| `RATE_LIMIT_BURST` | `20` | 순간 허용 요청 수 |
| `DB_QUERY_TIMEOUT` | `5s` | 쿼리 하나의 최대 실행 시간 |
| `SHUTDOWN_TIMEOUT` | `10s` | 종료 시 처리 중인 요청을 기다리는 시간 |

---

## 아키텍처

### 요청 흐름

```
HTTP 요청
  → middleware   Recover → Logging → CORS → RateLimit → (RequireAuth)
  → handlers     요청 파싱, 응답 작성
  → services     입력 검증, 권한 확인, 비즈니스 규칙
  → repository   SQL 실행 (인터페이스 + MySQL 구현)
  → MySQL
```

- 각 계층은 바로 아래 계층에만 의존합니다.
- 서비스는 repository **인터페이스**에만 의존하므로, 테스트에서는 가짜 구현을 끼워 DB 없이 검증할 수 있습니다.
- 의존성은 `main.go`에서 한 번에 조립합니다 (전역 DB 변수 없음).

### 에러 처리

- 서비스와 repository는 예상된 에러를 `apperr` 패키지로 만듭니다 (`Invalid`, `Unauthorized`, `Forbidden`, `NotFound`, `Conflict`).
- `httpx.WriteError`가 에러 종류를 HTTP 상태 코드로 바꾸고 메시지를 그대로 내려줍니다.
- 그 밖의 에러(DB 장애 등)는 서버 로그에만 남기고, 클라이언트에는 `"서버 내부 오류가 발생했습니다."`만 보냅니다.

### 디렉터리 구조

```
.
├── main.go               설정 로드, 의존성 조립, 서버 실행, graceful shutdown
├── router.go             URL 라우팅 (페이지 / 공개 API / 인증 API)
├── config/               환경 변수 → Config 구조체, 값 검증
├── db/
│   ├── connection.go     MySQL 커넥션 풀 생성
│   ├── schema.sql        테이블 정의
│   └── seed.sql          샘플 데이터
├── models/               도메인 모델, 요청/응답 DTO
├── repository/           DB 접근 계층 (인터페이스 + MySQL 구현)
├── services/             비즈니스 로직, 입력 검증, 권한 확인 (+ 단위 테스트)
├── handlers/             HTTP 핸들러, HTML 페이지 제공
├── middleware/           인증, 로깅, panic 복구, CORS, rate limit
├── auth/                 JWT 발급/검증, 비밀번호 해시
├── apperr/               에러 종류 정의
├── httpx/                JSON 응답 / 에러 / 요청 파싱 헬퍼
├── templates/            HTML 페이지 5개
├── Dockerfile            멀티 스테이지 빌드 (alpine, non-root 실행)
├── docker-compose.yml    MySQL + 앱
├── .env.example          환경 변수 예시
└── air.toml              개발용 자동 재빌드 설정
```

---

## 데이터베이스

전체 정의는 [`db/schema.sql`](db/schema.sql)에 있습니다.

```
users ──< user_rest >── restaurants ──< reservation
```

| 테이블 | 설명 | 주요 컬럼 / 제약 |
|---|---|---|
| `users` | 사장님 계정 | `email` UNIQUE, `password`(bcrypt 해시) |
| `restaurants` | 식당 | `name`, `address`, `description`, `phone_num`, `waiting_time`(분) |
| `user_rest` | 사장님 ↔ 식당 연결 | PK(`user_id`, `rest_id`), `rest_id` UNIQUE (식당당 사장님 1명) |
| `reservation` | 웨이팅 | `count`, `customer_phone_num`, `is_entry` (0 대기 / 1 입장 / 2 취소) |

- 사용자를 삭제하면 `user_rest`의 연결도 함께 삭제됩니다. 식당을 삭제하면 연결과 웨이팅이 함께 삭제됩니다 (`ON DELETE CASCADE`).
- 인덱스: `reservation(rest_id, is_entry)`, `reservation(rest_id, customer_phone_num, is_entry)`

---

## 인증과 권한

1. `POST /api/login`으로 로그인하면 `token`(JWT)을 받습니다.
2. 이후 보호된 API에는 헤더를 붙입니다.
   ```
   Authorization: Bearer <token>
   ```
3. 서버는 토큰에서 사용자 ID를 꺼내 사용합니다. URL로 다른 사람의 ID를 넘길 수 없습니다.

| 규칙 | 내용 |
|---|---|
| 본인만 | 내 정보 조회·수정·탈퇴는 토큰의 사용자 본인에게만 적용 |
| 사장님만 | 웨이팅 목록 조회와 상태 변경은 그 식당에 연결된 사장님만 가능 (아니면 403) |
| 식당당 1명 | 이미 사장님이 있는 식당에는 연결할 수 없음 (409) |
| 토큰 만료 | 기본 24시간. 만료되거나 위조된 토큰은 401 |

- 비밀번호는 bcrypt 해시로만 저장하며, 어떤 응답에도 포함되지 않습니다.
- 로그인 실패 시 "이메일이 없음"과 "비밀번호 틀림"을 구분하지 않고 같은 메시지를 보냅니다.

---

## API 명세

- 요청과 응답 본문은 모두 JSON입니다.
- 🔒 표시가 있는 API는 `Authorization: Bearer <token>` 헤더가 필요합니다.

### 요약

| 메서드 | 경로 | 인증 | 설명 | 성공 |
|---|---|---|---|---|
| GET | `/api/healthz` | | 서버 상태 확인 | 200 |
| POST | `/api/register` | | 회원가입 | 201 |
| POST | `/api/login` | | 로그인, 토큰 발급 | 200 |
| GET | `/api/user` | 🔒 | 내 정보 | 200 |
| PUT | `/api/user` | 🔒 | 내 정보 수정 | 200 |
| DELETE | `/api/user` | 🔒 | 회원 탈퇴 | 204 |
| POST | `/api/link` | 🔒 | 내 계정을 식당 사장님으로 등록 | 201 |
| GET | `/api/restaurant-list` | | 식당 목록 | 200 |
| GET | `/api/restaurant/{id}` | | 식당 상세 | 200 |
| POST | `/api/reservation` | | 웨이팅 등록 | 201 |
| GET | `/api/restaurant/{id}/reservations` | 🔒 | 내 식당 웨이팅 목록 | 200 |
| PUT | `/api/reservation/{id}/status` | 🔒 | 웨이팅 상태 변경 | 200 |

### 상태 코드

| 코드 | 의미 | 예시 |
|---|---|---|
| 400 | 입력 오류 | 잘못된 JSON, 잘못된 ID, 형식이 틀린 전화번호, 인원수 범위 초과 |
| 401 | 인증 필요 | 토큰 없음·만료·위조, 로그인 실패 |
| 403 | 권한 없음 | 남의 식당 웨이팅 조회·변경 |
| 404 | 없음 | 존재하지 않는 식당·예약·사용자·API |
| 405 | 메서드 불일치 | `PATCH /api/user` |
| 409 | 중복 | 이미 가입된 이메일, 이미 대기 중인 번호, 이미 사장님이 있는 식당 |
| 429 | 요청 과다 | rate limit 초과 (`Retry-After` 헤더 포함) |
| 500 | 서버 오류 | 내부 오류 (상세 내용은 서버 로그에만 기록) |

에러 응답 형식:
```json
{"error": "이미 이 식당에 대기 중인 예약이 있습니다."}
```

### 사용자

#### `POST /api/register` 회원가입

```json
{
  "email": "owner@example.com",
  "password": "password123",
  "name": "홍길동",
  "phone_num": "010-1234-5678"
}
```

- `email`: 필수. 소문자로 바꿔 저장합니다.
- `password`: 필수, 8~72자.
- `name`: 필수, 1~100자.
- `phone_num`: 선택. 하이픈을 제거해 저장합니다.

응답 `201`:
```json
{
  "id": 1,
  "email": "owner@example.com",
  "phone_num": "01012345678",
  "name": "홍길동",
  "rest_ids": [],
  "created_at": "2026-11-24T00:48:29Z",
  "updated_at": "2026-11-24T00:48:29Z"
}
```

#### `POST /api/login` 로그인

```json
{"email": "owner@example.com", "password": "password123"}
```

응답 `200`:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": { "id": 1, "email": "owner@example.com", "name": "홍길동", "rest_ids": [1], "...": "..." }
}
```

`user.rest_ids`는 이 사용자가 사장님으로 등록된 식당 ID 목록입니다.

#### `GET /api/user` 🔒 내 정보

응답은 회원가입 응답과 같은 형식입니다.

#### `PUT /api/user` 🔒 내 정보 수정

```json
{"email": "new@example.com", "name": "홍길동", "phone_num": "010-5555-6666"}
```

`email`과 `name`은 필수입니다. 다른 사람이 쓰는 이메일이면 409를 반환합니다.

#### `DELETE /api/user` 🔒 회원 탈퇴

응답 `204` (본문 없음). 식당 연결도 함께 삭제됩니다.

#### `POST /api/link` 🔒 식당 사장님 등록

```json
{"rest_id": 1}
```

응답 `201`:
```json
{"user_id": 1, "rest_id": 1, "created_at": "2026-11-24T00:48:29Z"}
```

없는 식당이면 404, 이미 사장님이 있는 식당이면 409를 반환합니다.

### 식당

#### `GET /api/restaurant-list` 식당 목록

응답 `200`:
```json
[
  {
    "id": 1,
    "name": "경희 국밥",
    "address": "서울 동대문구 경희대로 1",
    "description": "24시간 끓이는 진한 국밥집",
    "phone_num": "02-111-1111",
    "waiting_time": 15,
    "waiting_count": 3,
    "created_at": "...",
    "updated_at": "..."
  }
]
```

`waiting_count`는 현재 대기 중(`is_entry = 0`)인 팀 수입니다.

#### `GET /api/restaurant/{id}` 식당 상세

응답은 목록의 항목 하나와 같은 형식입니다.

### 웨이팅

#### `POST /api/reservation` 웨이팅 등록 (손님, 로그인 불필요)

```json
{"rest_id": 1, "count": 2, "customer_phone_num": "010-1234-5678"}
```

- `count`: 1~20명.
- `customer_phone_num`: 휴대폰 번호(010, 011, 016~019). 하이픈과 공백은 제거해 저장합니다.
- 같은 식당에 같은 번호로 이미 대기 중이면 409를 반환합니다. 다른 식당에는 등록할 수 있습니다.

응답 `201`:
```json
{"id": 1, "rest_id": 1, "count": 2, "customer_phone_num": "01012345678", "created_at": "...", "is_entry": 0}
```

#### `GET /api/restaurant/{id}/reservations` 🔒 웨이팅 목록 (사장님)

- 대기 중인 웨이팅이 먼저 온 순서로 위에 나옵니다.
- 처리된(입장·취소) 웨이팅은 그 아래에 최신순으로 나옵니다.

#### `PUT /api/reservation/{id}/status` 🔒 상태 변경 (사장님)

```json
{"is_entry": 1}
```

| 값 | 의미 |
|---|---|
| 0 | 대기 |
| 1 | 입장 |
| 2 | 취소 |

응답 `200`:
```json
{"id": 1, "is_entry": 1}
```

### curl 예시

```bash
curl -X POST localhost:8080/api/register -H 'Content-Type: application/json' -d '{"email":"owner@example.com","password":"password123","name":"홍길동"}'
```

```bash
TOKEN=$(curl -s -X POST localhost:8080/api/login -H 'Content-Type: application/json' -d '{"email":"owner@example.com","password":"password123"}' | jq -r .token)
```

```bash
curl -X POST localhost:8080/api/link -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' -d '{"rest_id":1}'
```

```bash
curl -X POST localhost:8080/api/reservation -H 'Content-Type: application/json' -d '{"rest_id":1,"count":2,"customer_phone_num":"010-1234-5678"}'
```

```bash
curl localhost:8080/api/restaurant/1/reservations -H "Authorization: Bearer $TOKEN"
```

```bash
curl -X PUT localhost:8080/api/reservation/1/status -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' -d '{"is_entry":1}'
```

---

## 화면

| 경로 | 파일 | 설명 |
|---|---|---|
| `/` | `templates/index.html` | 시작 화면 (식당 보기 / 사장님 로그인) |
| `/restaurant-list` | `templates/restaurant_list.html` | 식당 목록, 현재 대기 팀 수 |
| `/restaurant/{id}` | `templates/restaurant_detail.html` | 식당 상세, 웨이팅 등록 폼 |
| `/login` | `templates/login.html` | 사장님 로그인. 토큰을 `localStorage`에 저장한 뒤 내 식당 웨이팅 목록으로 이동 |
| `/reservation-list/{id}` | `templates/reservation_list.html` | 웨이팅 목록, 입장·취소 버튼. 토큰이 없거나 만료되면 로그인 화면으로 이동 |

---

## 테스트

```bash
go test ./...
```

| 패키지 | 검증 내용 |
|---|---|
| `services` | 웨이팅 입력 검증, 전화번호 정규화, 중복 웨이팅 409, 사장님 권한(403), 회원가입·로그인, 식당 미연결 사용자 로그인, 식당당 사장님 1명 |
| `auth` | 토큰 발급·검증, 다른 키로 서명된 토큰 거부, 만료 토큰 거부 |
| `middleware` | 인증 헤더 검사, panic 복구, rate limit, CORS preflight |
| `httpx` | 에러 종류 → 상태 코드 변환, 내부 에러 메시지 숨김 |

서비스 테스트는 repository 인터페이스의 가짜 구현(`services/fakes_test.go`)을 사용하므로 DB가 필요 없습니다.

---

## 알려진 한계와 다음 단계

**현재 한계**
- 같은 번호로 **동시에** 웨이팅을 등록하는 경쟁 상황은 막지 못합니다 (확인 후 삽입 방식).
- rate limit은 메모리에 저장하므로 서버 한 대에서만 정확합니다. 프록시 뒤에 배포하면 `X-Forwarded-For` 처리가 필요합니다.
- 회원가입 화면과 식당 등록 기능이 없습니다 (식당은 DB에 직접 넣어야 함).
- 식당 사장님 등록은 "먼저 등록한 사람" 방식이며, 관리자 승인 절차가 없습니다.

**다음 단계 후보**
- repository 통합 테스트 (testcontainers로 MySQL 띄우기), GitHub Actions CI
- OpenAPI(Swagger) 문서
- 웨이팅 상태 확장 (호출, 노쇼), 손님 셀프 조회·취소
- 사장님 화면 실시간 갱신 (WebSocket 또는 SSE), 입장 알림 (SMS, 알림톡)
- 식당 등록·수정, 검색, 페이지네이션
- 배포 (클라우드 + 관리형 MySQL), 구조화 로깅, 모니터링

---

## 기여자

- [GDG BE] 김도영
- [GDG BE] 이준호
