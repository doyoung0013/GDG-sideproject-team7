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
- Database: MySQL


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
backend-side-project/
├── main.go           # 메인 엔트리 포인트
├── go.mod            # Go 모듈 설정 파일
├── handlers/         # HTTP 핸들러 함수 모음
├── services/         # 비즈니스 로직 (서비스 계층)
├── models/           # 데이터 모델 정의
├── db/               # 데이터베이스 관련 로직
│   └── connection.go # DB 연결 설정
└── config/           # 설정 파일 및 환경 변수
```


## 실행 방법

```
- .env 파일 추가
```

서버는 8080 포트에서 실행됩니다.


## 기여자

- 김도영
- 이준호

