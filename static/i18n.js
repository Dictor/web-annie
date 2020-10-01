const i18nText = {
    ko: {
      message: {
        timeLeft: "남음",
        statusWait: "대기",
        statusProgress: "진행",
        statusComplete: "완료",
        statusFail: "실패",
        statusCancel: "취소",
        confirmDeleteTask: "{name} ({address}) 작업을 삭제합니까?",
        errorRefreshTask: "갱신 오류 : {msg}",
        promptTaskAddress: "추가할 동영상의 주소를 입력하세요.",
        infoAddSuccess: "추가 성공!",
        infoDeleteSuccess: "삭제 성공!",
        confirmDeleteCompletedTask: "완료된 모든 작업을 삭제합니까?",
        infoDeleteCompletedTask: "{count}개 삭제 성공!",
        error400: "유효하지 않은 입력 값",
        error500: "서버 내부 오류",
        errorUnknown: "예기치 못한 오류 = {status} {data}",
        errorRequest: "요청 중 오류", 
        errorGeneral: "실패 : {msg}"
      }
    },
    en: {
      message: {
        hello: 'こんにちは、世界'
      }
    }
  }

const i18n = new VueI18n({
    locale: navigator.language.split('-')[0],
    fallbackLocale: 'en',
    messages: i18nText,
});