# Personal AI Team Workspace

**AI 팀 워크플로**를 위한 구조화된 저장소다. 여러 AI 역할이 연구, 아키텍처, 구현 계획, 문서화에 협업하며, 소프트웨어 개발뿐 아니라 연구·설계·분석·작성을 아우르는 **구조화된 문제 해결**을 목표로 한다.

---

## 문서화 언어 (Documentation language)

이 저장소의 **설명용 문서는 한국어**를 기준으로 한다. **기술 식별자**는 영어를 유지한다.

- **한국어:** 설명, 목적, 규칙, 개요 등 모든 설명 문구.
- **영어 유지:** 도구명(ClusterLoader2, Prometheus, Grafana 등), metrics, SLI/SLO 용어, 명령어, 파일명, 설정 키.

적용 범위: `programs/`, `projects/`, `docs/` 아래 README 및 설명 문서. 기술 식별자는 번역하지 않는다.

---

## 개요 (Overview)

이 워크스페이스는 다음을 지원한다.

- **Research** — 도구 조사, 비교, 근거 수집
- **Architecture design** — 시스템·솔루션 설계
- **Technical analysis** — 타당성, 트레이드오프, 리스크
- **Implementation planning** — 작업, 마일스톤, 산출물
- **Documentation** — 보고서, spec, 지식 정리

여러 **AI 역할**이 순차(또는 병렬)로 동작하며, 각 역할은 명확한 prompt와 출력 위치를 가진다.

---

## 저장소 구조 (Repository Structure)

**ai-team-lab**은 다음 네 가지 영역으로 구성된다.

| 영역 | 역할 |
|------|------|
| **programs/** | 장기 작업 흐름(예: CAT, SRE Monitoring). 각 program은 여러 project를 포함하는 도메인이다. |
| **projects/** | 구체적인 문제 해결 단위. 각 project는 워크플로(manager → research → architecture → engineering → experiment → review → documentation)의 **실행 단위**이다. |
| **tasks/** | Task intake. 워크플로에 들어가기 전에 작업을 정의하는 곳. |
| **knowledge/** | 공유 학습. 완료된 project에서 뽑은 원칙, 도구, 패턴, 교훈. |

- **Program → 여러 Project를 포함.** Program은 도메인별로 관련 작업을 묶고, project는 `projects/` 아래에 두며 어느 program에 속하는지 선언한다.
- **Project → 워크플로로 실행.** Project는 task(또는 ad-hoc)에서 생성되며 **WORKFLOW.md**의 단계를 따른다.

```
ai-team-lab/
├── README.md                 # 이 파일
├── WORKFLOW.md               # AI 팀 내 작업 흐름
├── .cursor/rules/            # Cursor 운영 규칙 — AI 팀 동작 (아래 참고)
├── programs/                 # 장기 작업 흐름 (도메인)
│   ├── cat/                  # Cluster Acceptance Testing
│   │   └── README.md
│   └── sre-monitoring/       # SRE Monitoring
│       └── README.md
├── tasks/                    # Task intake — 워크플로 전 작업 정의
│   ├── README.md
│   ├── intake-template.md    # 새 task 표준 형식
│   └── example-kubernetes-cat.md
├── scripts/                  # 워크플로 실행·헬퍼
│   ├── README.md
│   └── run_workflow.py       # task 파일로 project 초기화
├── prompts/                  # 각 AI 역할용 prompt
│   ├── manager.md
│   ├── researcher.md
│   ├── architect.md
│   ├── engineer.md
│   ├── reviewer.md
│   └── writer.md
├── templates/                # 재사용 문서 템플릿
│   ├── research.md
│   ├── architecture.md
│   └── final-report.md
├── knowledge/                # 지식 메모리 — project에서 나온 재사용 학습
│   ├── README.md
│   ├── principles/           # 설계 원칙, 결정 규칙
│   ├── tools/                # 도구 평가·사용 노트
│   ├── patterns/             # 재사용 솔루션 패턴
│   └── lessons/              # 교훈
├── projects/                 # project당 한 폴더; 단계별 출력
│   ├── README.md
│   └── _sample/              # 예제 project
│       ├── README.md
│       ├── 01-manager/
│       ├── 02-research/
│       ├── 03-architecture/
│       ├── 04-engineering/
│       ├── 05-review/
│       └── 06-documentation/
```

---

## 역할 (Roles)

| Role | 초점 | 대표 산출물 |
|------|------|-------------|
| **Manager** | Task 분석, 분해, 우선순위 | Brief, 작업 분해, handoff |
| **Researcher** | 도구 조사, 비교, 근거 | Research notes, 비교표 |
| **Architect** | 시스템/솔루션 설계, 제약 | Architecture doc, 다이어그램 |
| **Engineer** | 구현 계획, 코드, 설정 | Spec, 코드, runbook |
| **Reviewer** | 검증, 비판, 공백 식별 | Review notes, 체크리스트 |
| **Writer** | 문서, 보고서, 요약 | Final report, docs, README |

---

## 새 AI 팀 project 실행 (Running a new AI team project)

1. **Task 생성:** `tasks/intake-template.md`를 사용해 `tasks/`에 작성 (예: `tasks/my-task.md`).
2. **Project 초기화:** 워크플로 러너로 project 폴더와 단계 디렉터리를 만들고 템플릿을 복사한다.

   ```bash
   python scripts/run_workflow.py tasks/my-task.md
   ```

   이렇게 하면 `projects/my-task/`가 생기고 `01-manager`부터 `06-documentation`까지 구성되며, task 파일이 `task-intake.md`로 복사되고 research, architecture, final-report 템플릿이 해당 단계에 들어간다.
3. **각 단계 실행:** `prompts/`의 prompt와 단계 폴더의 템플릿을 사용한다. Manager(입력: `task-intake.md`)부터 시작한 뒤 Researcher, Architect, Engineer, Reviewer, Writer, 마지막으로 **Knowledge Extraction** 순서. 전체 과정은 **WORKFLOW.md** 참고.
4. **Project 완료 후** **Knowledge Extraction** 실행: `templates/knowledge-extraction.md`를 채우고(예: `07-knowledge-extraction/` 또는 `08-knowledge-extraction/`), `knowledge/principles/`, `knowledge/tools/`, `knowledge/patterns/`, `knowledge/lessons/`에 새로 만들거나 갱신한다. 이렇게 해서 project가 재사용 가능한 학습으로 남고 AI 팀이 점진적으로 개선된다. 각 디렉터리 용도는 `knowledge/README.md` 참고.

스크립트는 **워크플로 초기화만** 수행한다. AI를 실행하거나 단계를 자동화하지 않으며, 워크플로를 단계별로 실행할 수 있도록 준비만 한다.

---

## Quick Start

1. **Task(intake) 작성:** `tasks/intake-template.md`를 `tasks/`에 새 파일로 복사 (예: `tasks/my-task-name.md`). **Task Title**, **Problem Description**, **Goal**, **Scope**, **Expected Deliverables**, **Constraints**, **Priority**, **Additional Context**를 채운다. 워크플로 진입 전 작업 정의의 표준 방식이다. `tasks/README.md`와 예제 `tasks/example-kubernetes-cat.md` 참고.
2. **Project 생성:** `python scripts/run_workflow.py tasks/<task-name>.md`로 project와 단계 폴더를 만든다 (위 **새 AI 팀 project 실행** 참고). 또는 `projects/<project-name>/`와 단계 폴더를 수동으로 만든다.
3. **워크플로 실행:** **WORKFLOW.md**를 따른다. Manager에 task intake를 입력한 뒤 Researcher, Architect, Engineer, Reviewer, Writer를 해당 폴더·템플릿으로 실행한다.
4. **Project마다 Knowledge Extraction:** Knowledge Extraction 단계와 `templates/knowledge-extraction.md`를 사용하고, 결과를 `knowledge/principles/`, `knowledge/tools/`, `knowledge/patterns/`, `knowledge/lessons/`에 저장한다. **WORKFLOW.md**와 `knowledge/README.md` 참고.
5. **템플릿 사용:** 필요 시 `templates/`에서 project 단계 폴더(예: `02-research/`)로 복사한다.

---

## Cursor와 AI 팀 워크플로

이 저장소에서 **Cursor**로 작업할 때, 어시스턴트는 `.cursor/rules/`의 **운영 규칙**을 따르는 것으로 간주된다. 해당 규칙에 따라 Cursor가 이 AI 팀 워크스페이스의 기본 운영자처럼 동작한다.

- 가능하면 `tasks/`의 task intake에서 시작
- **WORKFLOW.md**의 순서를 따르고, 요청이 없으면 단계를 건너뛰지 않음
- 각 역할(Manager, Researcher, Architect 등)의 책임 범위 유지
- 산출물은 `projects/<project-name>/0X-<phase>/`에 저장하고 템플릿 사용
- Project 완료 후 Knowledge Extraction 및 `knowledge/` 갱신 고려
- 복잡함보다 명확성, 구조, 재사용 우선

동작을 바꾸려면 `.cursor/rules/`의 규칙을 검토·수정하면 된다.

---

## 설계 원칙 (Design Principles)

- **Simple** — 최소 구조, 탐색·확장이 쉬움.
- **Readable** — 명확한 이름, 짧은 문서, Markdown 위주.
- **Extensible** — 역할, 템플릿, 단계를 추가해도 워크플로가 깨지지 않음.

전체 task 흐름은 **WORKFLOW.md**, 각 역할 지침은 **prompts/**를 참고한다.
