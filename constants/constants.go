package constants

type ImageTypeIds int32

const (
	ImageTypeIdDefault           ImageTypeIds = 100
	ImageTypeIdUserAvatar        ImageTypeIds = 200
	ImageTypeIdWorldHeader       ImageTypeIds = 300
	ImageTypeIdWorldAvatar       ImageTypeIds = 400
	ImageTypeIdLocationImage     ImageTypeIds = 500
	ImageTypeIdRaceImage         ImageTypeIds = 600
	ImageTypeIdItemImage         ImageTypeIds = 700
	ImageTypeIdSkillImage        ImageTypeIds = 800
	ImageTypeIdCharacterPortrait ImageTypeIds = 900
	ImageTypeIdMapImage          ImageTypeIds = 1000
	ImageTypeIdBackgroundImage   ImageTypeIds = 1100
	ImageTypeIdWorldThumbnail    ImageTypeIds = 1200
	ImageTypeIdMenuHeader        ImageTypeIds = 1300
)

const UniversalWorldName = "Universal World"
const UniversalSystemName = "Universal System"

const GrpcCookieHeader = "grpc-gateway-cookie"
const CookieName = "access_token"
