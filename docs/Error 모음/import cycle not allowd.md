## import Cycle
![import cycle error](../img/import cycle.png)

![import cycle error](../img/import cycle error ppt.png)
- Core 모듈에서 Extension 을 불러오고, Extension 모듈에서 Core 모듈을 불러오면
- 순환 에러가 발생한다.
- 이를 해결하기 위해서 Type 으로 정의해놓은 상수 변수들을 Type 패키지로 따로 빼놓았다.

