# Restaurant Reservation Service

Go로 구현된 레스토랑 예약 관리 시스템입니다.

## 주요 기능

- 레스토랑 목록 조회 및 상세 정보 확인
- 사용자 회원가입/로그인
- 레스토랑 예약 생성 및 관리
- 예약 상태 업데이트
- 레스토랑-사용자 연동

## 기술 스택

- Backend: Go
- Router: Gorilla Mux
- Database: SQL

## API 엔드포인트

### 레스토랑
```
GET    /api/restaurant-list         # 레스토랑 목록 조회
GET    /api/restaurant/{id}         # 레스토랑 상세 정보
POST   /api/reservation            # 예약 생성
PUT    /api/reservation/{id}/status # 예약 상태 수정
```

### 사용자
```
POST   /api/register               # 회원가입
POST   /api/login                 # 로그인
GET    /api/user/{id}             # 사용자 정보 조회
PUT    /api/user                  # 사용자 정보 수정
DELETE /api/user/{id}             # 회원 탈퇴
```

### 예약
```
GET    /api/user/reservation-list/{id}  # 예약 목록 조회
POST   /api/link                        # 레스토랑-사용자 연동
```

## 프로젝트 구조

```
├── config/          # 설정 파일
├── db/              # 데이터베이스 연결
├── handlers/        # HTTP 핸들러
├── services/        # 비즈니스 로직
└── main.go         # 애플리케이션 진입점
```

## 실행 방법

```bash
git clone https://github.com/GDG-KHU-Side/backend-side-project.git
cd backend-side-project
go mod download
go run main.go
```

서버는 8080 포트에서 실행됩니다.

## 기여자

- [기여자1 이름]
- [기여자2 이름]

## 라이선스

MIT License
