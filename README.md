# GDG Backend Side Project

Go로 구현된 레스토랑 예약 관리 시스템입니다.


## 주요 기능

- 사용자 회원가입/로그인
- 레스토랑 목록 조회 및 상세 정보 확인
- 레스토랑 예약 생성 및 관리
- 예약 상태 업데이트
- 레스토랑-사용자 연동


## 기술 스택

- Backend: Go 1.23, Gorilla Mux
- Database: MySQL 8
- 인증: JWT (HS256), bcrypt
- 기타: golang.org/x/time/rate (요청 속도 제한), Docker Compose


## 실행 방법

### Docker Compose (권장)
```bash
cp .env.example .env   # JWT_SECRET 등 값 수정
docker compose up --build
```
MySQL이 뜨면서 `db/schema.sql`, `db/seed.sql`(샘플 식당 3개)이 자동 적용되고, 서버는 http://localhost:8080 에서 실행됩니다.
스키마를 다시 적용하려면 `docker compose down -v` 로 볼륨을 지운 뒤 다시 띄우세요.

### 로컬 실행
```bash
cp .env.example .env              # DB 접속 정보, JWT_SECRET 설정
mysql -u root -p gdgdb < db/schema.sql
mysql -u root -p gdgdb < db/seed.sql   # 선택
go run .                          # 또는 air (자동 재빌드)
```

### 테스트
```bash
go test ./...
```


## API

모든 에러 응답은 `{"error": "메시지"}` 형식입니다.
🔒 표시는 `Authorization: Bearer <token>` 헤더가 필요합니다. 토큰은 `/api/login` 응답으로 받습니다.

### 사용자
| 메서드 | 경로 | 설명 | 성공 |
|---|---|---|---|
| POST | `/api/register` | 회원가입 `{email, password(8자 이상), name, phone_num?}` | 201 |
| POST | `/api/login` | 로그인 → `{token, user}` (`user.rest_ids`: 연결된 식당) | 200 |
| GET 🔒 | `/api/user` | 내 정보 | 200 |
| PUT 🔒 | `/api/user` | 내 정보 수정 `{email, name, phone_num?}` | 200 |
| DELETE 🔒 | `/api/user` | 회원 탈퇴 | 204 |
| POST 🔒 | `/api/link` | 내 계정을 식당 사장님으로 등록 `{rest_id}` (식당당 1명, 이미 있으면 409) | 201 |

### 식당 / 웨이팅
| 메서드 | 경로 | 설명 | 성공 |
|---|---|---|---|
| GET | `/api/restaurant-list` | 식당 목록 (`waiting_count`: 현재 대기 팀 수) | 200 |
| GET | `/api/restaurant/{id}` | 식당 상세 | 200 |
| POST | `/api/reservation` | 웨이팅 등록 `{rest_id, count(1~20), customer_phone_num}` (같은 번호로 대기 중이면 409) | 201 |
| GET 🔒 | `/api/restaurant/{id}/reservations` | 내 식당 웨이팅 목록 (사장님만) | 200 |
| PUT 🔒 | `/api/reservation/{id}/status` | 상태 변경 `{is_entry: 0 대기 / 1 입장 / 2 취소}` (사장님만) | 200 |
| GET | `/api/healthz` | 상태 확인 | 200 |

주요 에러 코드: 400(입력 오류), 401(로그인 필요), 403(권한 없음), 404(없음), 409(중복), 429(요청 과다)


## 프로젝트 구조

```
main.go          설정 로드, 의존성 조립, 서버 실행/graceful shutdown
router.go        라우팅
config/          환경변수 설정 (Config)
db/              DB 연결, schema.sql, seed.sql
models/          도메인 모델과 요청/응답 DTO
repository/      DB 접근 (인터페이스 + MySQL 구현)
services/        비즈니스 로직, 입력 검증, 권한 확인
handlers/        HTTP 핸들러
middleware/      인증(JWT), 로깅, panic 복구, CORS, rate limit
auth/            JWT 발급/검증, 비밀번호 해시
apperr/          에러 종류 (→ HTTP 상태 코드)
httpx/           JSON 응답/에러 헬퍼
templates/       HTML 페이지
```

요청 흐름: `middleware → handlers → services → repository → MySQL`


## 기여자

- [GDG BE] 김도영
- [GDG BE] 이준호

