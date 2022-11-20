package main

type Vector2 struct {
	x float32
	y float32
}

type Vector3 struct {
	x float32
	y float32
	z float32
}

type Vector4 struct {
	x float32
	y float32
	z float32
	w float32
}

type Quaternion = Vector4

type Matrix struct {
	m0  float32
	m4  float32
	m8  float32
	m12 float32
	m1  float32
	m5  float32
	m9  float32
	m13 float32
	m2  float32
	m6  float32
	m10 float32
	m14 float32
	m3  float32
	m7  float32
	m11 float32
	m15 float32
}

type Color struct {
	r uint8
	g uint8
	b uint8
	a uint8
}

type Rectangle struct {
	x      float32
	y      float32
	width  float32
	height float32
}

type Image struct {
	data    interface{}
	width   int32
	height  int32
	mipmaps int32
	format  int32
}

type Texture struct {
	id      uint32
	width   int32
	height  int32
	mipmaps int32
	format  int32
}

type Texture2D = Texture

type TextureCubemap = Texture

type RenderTexture struct {
	id      uint32
	texture Texture
	depth   Texture
}

type RenderTexture2D = RenderTexture

type NPatchInfo struct {
	source Rectangle
	left   int32
	top    int32
	right  int32
	bottom int32
	layout int32
}

type GlyphInfo struct {
	value    int32
	offsetX  int32
	offsetY  int32
	advanceX int32
	image    Image
}

type Font struct {
	baseSize     int32
	glyphCount   int32
	glyphPadding int32
	texture      Texture2D
	recs         []Rectangle
	glyphs       []GlyphInfo
}

type Camera3D struct {
	position   Vector3
	target     Vector3
	up         Vector3
	fovy       float32
	projection int32
}

type Camera = Camera3D

type Camera2D struct {
	offset   Vector2
	target   Vector2
	rotation float32
	zoom     float32
}

type Mesh struct {
	vertexCount   int32
	triangleCount int32
	vertices      []float32
	texcoords     []float32
	texcoords2    []float32
	normals       []float32
	tangents      []float32
	colors        []uint8
	indices       []uint16
	animVertices  []float32
	animNormals   []float32
	boneIds       []uint8
	boneWeights   []float32
	vaoId         uint32
	vboId         []uint32
}

type Shader struct {
	id   uint32
	locs []int32
}

type MaterialMap struct {
	texture Texture2D
	color   Color
	value   float32
}

type Material struct {
	shader Shader
	maps   []MaterialMap
	params [4]float32
}

type Transform struct {
	translation Vector3
	rotation    Quaternion
	scale       Vector3
}

type BoneInfo struct {
	name   [32]byte
	parent int32
}

type Model struct {
	transform     Matrix
	meshCount     int32
	materialCount int32
	meshes        []Mesh
	materials     []Material
	meshMaterial  []int32
	boneCount     int32
	bones         []BoneInfo
	bindPose      []Transform
}

type ModelAnimation struct {
	boneCount  int32
	frameCount int32
	bones      []BoneInfo
	framePoses [][]Transform
}

type Ray struct {
	position  Vector3
	direction Vector3
}

type RayCollision struct {
	hit      int32
	distance float32
	point    Vector3
	normal   Vector3
}

type BoundingBox struct {
	min Vector3
	max Vector3
}

type Wave struct {
	frameCount uint32
	sampleRate uint32
	sampleSize uint32
	channels   uint32
	data       interface{}
}

type AudioStream struct {
	buffer     []rAudioBuffer
	processor  []rAudioProcessor
	sampleRate uint32
	sampleSize uint32
	channels   uint32
}

type Sound struct {
	stream     AudioStream
	frameCount uint32
}

type Music struct {
	stream     AudioStream
	frameCount uint32
	looping    int32
	ctxType    int32
	ctxData    interface{}
}

type VrDeviceInfo struct {
	hResolution            int32
	vResolution            int32
	hScreenSize            float32
	vScreenSize            float32
	vScreenCenter          float32
	eyeToScreenDistance    float32
	lensSeparationDistance float32
	interpupillaryDistance float32
	lensDistortionValues   [4]float32
	chromaAbCorrection     [4]float32
}

type VrStereoConfig struct {
	projection        [2]Matrix
	viewOffset        [2]Matrix
	leftLensCenter    [2]float32
	rightLensCenter   [2]float32
	leftScreenCenter  [2]float32
	rightScreenCenter [2]float32
	scale             [2]float32
	scaleIn           [2]float32
}

type FilePathList struct {
	capacity uint32
	count    uint32
	paths    [][]byte
}

const (
	FLAG_VSYNC_HINT               int32 = 64
	FLAG_FULLSCREEN_MODE                = 2
	FLAG_WINDOW_RESIZABLE               = 4
	FLAG_WINDOW_UNDECORATED             = 8
	FLAG_WINDOW_HIDDEN                  = 128
	FLAG_WINDOW_MINIMIZED               = 512
	FLAG_WINDOW_MAXIMIZED               = 1024
	FLAG_WINDOW_UNFOCUSED               = 2048
	FLAG_WINDOW_TOPMOST                 = 4096
	FLAG_WINDOW_ALWAYS_RUN              = 256
	FLAG_WINDOW_TRANSPARENT             = 16
	FLAG_WINDOW_HIGHDPI                 = 8192
	FLAG_WINDOW_MOUSE_PASSTHROUGH       = 16384
	FLAG_MSAA_4X_HINT                   = 32
	FLAG_INTERLACED_HINT                = 65536
)

type ConfigFlags = int32

const (
	LOG_ALL     int32 = 0
	LOG_TRACE         = 1
	LOG_DEBUG         = 2
	LOG_INFO          = 3
	LOG_WARNING       = 4
	LOG_ERROR         = 5
	LOG_FATAL         = 6
	LOG_NONE          = 7
)

type TraceLogLevel = int32

const (
	KEY_NULL          int32 = 0
	KEY_APOSTROPHE          = 39
	KEY_COMMA               = 44
	KEY_MINUS               = 45
	KEY_PERIOD              = 46
	KEY_SLASH               = 47
	KEY_ZERO                = 48
	KEY_ONE                 = 49
	KEY_TWO                 = 50
	KEY_THREE               = 51
	KEY_FOUR                = 52
	KEY_FIVE                = 53
	KEY_SIX                 = 54
	KEY_SEVEN               = 55
	KEY_EIGHT               = 56
	KEY_NINE                = 57
	KEY_SEMICOLON           = 59
	KEY_EQUAL               = 61
	KEY_A                   = 65
	KEY_B                   = 66
	KEY_C                   = 67
	KEY_D                   = 68
	KEY_E                   = 69
	KEY_F                   = 70
	KEY_G                   = 71
	KEY_H                   = 72
	KEY_I                   = 73
	KEY_J                   = 74
	KEY_K                   = 75
	KEY_L                   = 76
	KEY_M                   = 77
	KEY_N                   = 78
	KEY_O                   = 79
	KEY_P                   = 80
	KEY_Q                   = 81
	KEY_R                   = 82
	KEY_S                   = 83
	KEY_T                   = 84
	KEY_U                   = 85
	KEY_V                   = 86
	KEY_W                   = 87
	KEY_X                   = 88
	KEY_Y                   = 89
	KEY_Z                   = 90
	KEY_LEFT_BRACKET        = 91
	KEY_BACKSLASH           = 92
	KEY_RIGHT_BRACKET       = 93
	KEY_GRAVE               = 96
	KEY_SPACE               = 32
	KEY_ESCAPE              = 256
	KEY_ENTER               = 257
	KEY_TAB                 = 258
	KEY_BACKSPACE           = 259
	KEY_INSERT              = 260
	KEY_DELETE              = 261
	KEY_RIGHT               = 262
	KEY_LEFT                = 263
	KEY_DOWN                = 264
	KEY_UP                  = 265
	KEY_PAGE_UP             = 266
	KEY_PAGE_DOWN           = 267
	KEY_HOME                = 268
	KEY_END                 = 269
	KEY_CAPS_LOCK           = 280
	KEY_SCROLL_LOCK         = 281
	KEY_NUM_LOCK            = 282
	KEY_PRINT_SCREEN        = 283
	KEY_PAUSE               = 284
	KEY_F1                  = 290
	KEY_F2                  = 291
	KEY_F3                  = 292
	KEY_F4                  = 293
	KEY_F5                  = 294
	KEY_F6                  = 295
	KEY_F7                  = 296
	KEY_F8                  = 297
	KEY_F9                  = 298
	KEY_F10                 = 299
	KEY_F11                 = 300
	KEY_F12                 = 301
	KEY_LEFT_SHIFT          = 340
	KEY_LEFT_CONTROL        = 341
	KEY_LEFT_ALT            = 342
	KEY_LEFT_SUPER          = 343
	KEY_RIGHT_SHIFT         = 344
	KEY_RIGHT_CONTROL       = 345
	KEY_RIGHT_ALT           = 346
	KEY_RIGHT_SUPER         = 347
	KEY_KB_MENU             = 348
	KEY_KP_0                = 320
	KEY_KP_1                = 321
	KEY_KP_2                = 322
	KEY_KP_3                = 323
	KEY_KP_4                = 324
	KEY_KP_5                = 325
	KEY_KP_6                = 326
	KEY_KP_7                = 327
	KEY_KP_8                = 328
	KEY_KP_9                = 329
	KEY_KP_DECIMAL          = 330
	KEY_KP_DIVIDE           = 331
	KEY_KP_MULTIPLY         = 332
	KEY_KP_SUBTRACT         = 333
	KEY_KP_ADD              = 334
	KEY_KP_ENTER            = 335
	KEY_KP_EQUAL            = 336
	KEY_BACK                = 4
	KEY_MENU                = 82
	KEY_VOLUME_UP           = 24
	KEY_VOLUME_DOWN         = 25
)

type KeyboardKey = int32

const (
	MOUSE_BUTTON_LEFT    int32 = 0
	MOUSE_BUTTON_RIGHT         = 1
	MOUSE_BUTTON_MIDDLE        = 2
	MOUSE_BUTTON_SIDE          = 3
	MOUSE_BUTTON_EXTRA         = 4
	MOUSE_BUTTON_FORWARD       = 5
	MOUSE_BUTTON_BACK          = 6
)

type MouseButton = int32

const (
	MOUSE_CURSOR_DEFAULT       int32 = 0
	MOUSE_CURSOR_ARROW               = 1
	MOUSE_CURSOR_IBEAM               = 2
	MOUSE_CURSOR_CROSSHAIR           = 3
	MOUSE_CURSOR_POINTING_HAND       = 4
	MOUSE_CURSOR_RESIZE_EW           = 5
	MOUSE_CURSOR_RESIZE_NS           = 6
	MOUSE_CURSOR_RESIZE_NWSE         = 7
	MOUSE_CURSOR_RESIZE_NESW         = 8
	MOUSE_CURSOR_RESIZE_ALL          = 9
	MOUSE_CURSOR_NOT_ALLOWED         = 10
)

type MouseCursor = int32

const (
	GAMEPAD_BUTTON_UNKNOWN          int32 = 0
	GAMEPAD_BUTTON_LEFT_FACE_UP           = 1
	GAMEPAD_BUTTON_LEFT_FACE_RIGHT        = 2
	GAMEPAD_BUTTON_LEFT_FACE_DOWN         = 3
	GAMEPAD_BUTTON_LEFT_FACE_LEFT         = 4
	GAMEPAD_BUTTON_RIGHT_FACE_UP          = 5
	GAMEPAD_BUTTON_RIGHT_FACE_RIGHT       = 6
	GAMEPAD_BUTTON_RIGHT_FACE_DOWN        = 7
	GAMEPAD_BUTTON_RIGHT_FACE_LEFT        = 8
	GAMEPAD_BUTTON_LEFT_TRIGGER_1         = 9
	GAMEPAD_BUTTON_LEFT_TRIGGER_2         = 10
	GAMEPAD_BUTTON_RIGHT_TRIGGER_1        = 11
	GAMEPAD_BUTTON_RIGHT_TRIGGER_2        = 12
	GAMEPAD_BUTTON_MIDDLE_LEFT            = 13
	GAMEPAD_BUTTON_MIDDLE                 = 14
	GAMEPAD_BUTTON_MIDDLE_RIGHT           = 15
	GAMEPAD_BUTTON_LEFT_THUMB             = 16
	GAMEPAD_BUTTON_RIGHT_THUMB            = 17
)

type GamepadButton = int32

const (
	GAMEPAD_AXIS_LEFT_X        int32 = 0
	GAMEPAD_AXIS_LEFT_Y              = 1
	GAMEPAD_AXIS_RIGHT_X             = 2
	GAMEPAD_AXIS_RIGHT_Y             = 3
	GAMEPAD_AXIS_LEFT_TRIGGER        = 4
	GAMEPAD_AXIS_RIGHT_TRIGGER       = 5
)

type GamepadAxis = int32

const (
	MATERIAL_MAP_ALBEDO     int32 = 0
	MATERIAL_MAP_METALNESS        = 1
	MATERIAL_MAP_NORMAL           = 2
	MATERIAL_MAP_ROUGHNESS        = 3
	MATERIAL_MAP_OCCLUSION        = 4
	MATERIAL_MAP_EMISSION         = 5
	MATERIAL_MAP_HEIGHT           = 6
	MATERIAL_MAP_CUBEMAP          = 7
	MATERIAL_MAP_IRRADIANCE       = 8
	MATERIAL_MAP_PREFILTER        = 9
	MATERIAL_MAP_BRDF             = 10
)

type MaterialMapIndex = int32

const (
	SHADER_LOC_VERTEX_POSITION   int32 = 0
	SHADER_LOC_VERTEX_TEXCOORD01       = 1
	SHADER_LOC_VERTEX_TEXCOORD02       = 2
	SHADER_LOC_VERTEX_NORMAL           = 3
	SHADER_LOC_VERTEX_TANGENT          = 4
	SHADER_LOC_VERTEX_COLOR            = 5
	SHADER_LOC_MATRIX_MVP              = 6
	SHADER_LOC_MATRIX_VIEW             = 7
	SHADER_LOC_MATRIX_PROJECTION       = 8
	SHADER_LOC_MATRIX_MODEL            = 9
	SHADER_LOC_MATRIX_NORMAL           = 10
	SHADER_LOC_VECTOR_VIEW             = 11
	SHADER_LOC_COLOR_DIFFUSE           = 12
	SHADER_LOC_COLOR_SPECULAR          = 13
	SHADER_LOC_COLOR_AMBIENT           = 14
	SHADER_LOC_MAP_ALBEDO              = 15
	SHADER_LOC_MAP_METALNESS           = 16
	SHADER_LOC_MAP_NORMAL              = 17
	SHADER_LOC_MAP_ROUGHNESS           = 18
	SHADER_LOC_MAP_OCCLUSION           = 19
	SHADER_LOC_MAP_EMISSION            = 20
	SHADER_LOC_MAP_HEIGHT              = 21
	SHADER_LOC_MAP_CUBEMAP             = 22
	SHADER_LOC_MAP_IRRADIANCE          = 23
	SHADER_LOC_MAP_PREFILTER           = 24
	SHADER_LOC_MAP_BRDF                = 25
)

type ShaderLocationIndex = int32

const (
	SHADER_UNIFORM_FLOAT     int32 = 0
	SHADER_UNIFORM_VEC2            = 1
	SHADER_UNIFORM_VEC3            = 2
	SHADER_UNIFORM_VEC4            = 3
	SHADER_UNIFORM_INT             = 4
	SHADER_UNIFORM_IVEC2           = 5
	SHADER_UNIFORM_IVEC3           = 6
	SHADER_UNIFORM_IVEC4           = 7
	SHADER_UNIFORM_SAMPLER2D       = 8
)

type ShaderUniformDataType = int32

const (
	SHADER_ATTRIB_FLOAT int32 = 0
	SHADER_ATTRIB_VEC2        = 1
	SHADER_ATTRIB_VEC3        = 2
	SHADER_ATTRIB_VEC4        = 3
)

type ShaderAttributeDataType = int32

const (
	PIXELFORMAT_UNCOMPRESSED_GRAYSCALE    int32 = 1
	PIXELFORMAT_UNCOMPRESSED_GRAY_ALPHA         = 2
	PIXELFORMAT_UNCOMPRESSED_R5G6B5             = 3
	PIXELFORMAT_UNCOMPRESSED_R8G8B8             = 4
	PIXELFORMAT_UNCOMPRESSED_R5G5B5A1           = 5
	PIXELFORMAT_UNCOMPRESSED_R4G4B4A4           = 6
	PIXELFORMAT_UNCOMPRESSED_R8G8B8A8           = 7
	PIXELFORMAT_UNCOMPRESSED_R32                = 8
	PIXELFORMAT_UNCOMPRESSED_R32G32B32          = 9
	PIXELFORMAT_UNCOMPRESSED_R32G32B32A32       = 10
	PIXELFORMAT_COMPRESSED_DXT1_RGB             = 11
	PIXELFORMAT_COMPRESSED_DXT1_RGBA            = 12
	PIXELFORMAT_COMPRESSED_DXT3_RGBA            = 13
	PIXELFORMAT_COMPRESSED_DXT5_RGBA            = 14
	PIXELFORMAT_COMPRESSED_ETC1_RGB             = 15
	PIXELFORMAT_COMPRESSED_ETC2_RGB             = 16
	PIXELFORMAT_COMPRESSED_ETC2_EAC_RGBA        = 17
	PIXELFORMAT_COMPRESSED_PVRT_RGB             = 18
	PIXELFORMAT_COMPRESSED_PVRT_RGBA            = 19
	PIXELFORMAT_COMPRESSED_ASTC_4x4_RGBA        = 20
	PIXELFORMAT_COMPRESSED_ASTC_8x8_RGBA        = 21
)

type PixelFormat = int32

const (
	TEXTURE_FILTER_POINT           int32 = 0
	TEXTURE_FILTER_BILINEAR              = 1
	TEXTURE_FILTER_TRILINEAR             = 2
	TEXTURE_FILTER_ANISOTROPIC_4X        = 3
	TEXTURE_FILTER_ANISOTROPIC_8X        = 4
	TEXTURE_FILTER_ANISOTROPIC_16X       = 5
)

type TextureFilter = int32

const (
	TEXTURE_WRAP_REPEAT        int32 = 0
	TEXTURE_WRAP_CLAMP               = 1
	TEXTURE_WRAP_MIRROR_REPEAT       = 2
	TEXTURE_WRAP_MIRROR_CLAMP        = 3
)

type TextureWrap = int32

const (
	CUBEMAP_LAYOUT_AUTO_DETECT         int32 = 0
	CUBEMAP_LAYOUT_LINE_VERTICAL             = 1
	CUBEMAP_LAYOUT_LINE_HORIZONTAL           = 2
	CUBEMAP_LAYOUT_CROSS_THREE_BY_FOUR       = 3
	CUBEMAP_LAYOUT_CROSS_FOUR_BY_THREE       = 4
	CUBEMAP_LAYOUT_PANORAMA                  = 5
)

type CubemapLayout = int32

const (
	FONT_DEFAULT int32 = 0
	FONT_BITMAP        = 1
	FONT_SDF           = 2
)

type FontType = int32

const (
	BLEND_ALPHA             int32 = 0
	BLEND_ADDITIVE                = 1
	BLEND_MULTIPLIED              = 2
	BLEND_ADD_COLORS              = 3
	BLEND_SUBTRACT_COLORS         = 4
	BLEND_ALPHA_PREMULTIPLY       = 5
	BLEND_CUSTOM                  = 6
	BLEND_CUSTOM_SEPARATE         = 7
)

type BlendMode = int32

const (
	GESTURE_NONE        int32 = 0
	GESTURE_TAP               = 1
	GESTURE_DOUBLETAP         = 2
	GESTURE_HOLD              = 4
	GESTURE_DRAG              = 8
	GESTURE_SWIPE_RIGHT       = 16
	GESTURE_SWIPE_LEFT        = 32
	GESTURE_SWIPE_UP          = 64
	GESTURE_SWIPE_DOWN        = 128
	GESTURE_PINCH_IN          = 256
	GESTURE_PINCH_OUT         = 512
)

type Gesture = int32

const (
	CAMERA_CUSTOM       int32 = 0
	CAMERA_FREE               = 1
	CAMERA_ORBITAL            = 2
	CAMERA_FIRST_PERSON       = 3
	CAMERA_THIRD_PERSON       = 4
)

type CameraMode = int32

const (
	CAMERA_PERSPECTIVE  int32 = 0
	CAMERA_ORTHOGRAPHIC       = 1
)

type CameraProjection = int32

const (
	NPATCH_NINE_PATCH             int32 = 0
	NPATCH_THREE_PATCH_VERTICAL         = 1
	NPATCH_THREE_PATCH_HORIZONTAL       = 2
)

type NPatchLayout = int32

type TraceLogCallback = *va_list

type LoadFileDataCallback = func([]byte, []uint32) []uint8

type SaveFileDataCallback = func([]byte, interface{}, uint32) int32

type LoadFileTextCallback = func([]byte) []byte

type SaveFileTextCallback = func([]byte, []byte) int32

type AudioCallback = func(interface{}, uint32)

type GuiStyleProp struct {
	controlId     uint16
	propertyId    uint16
	propertyValue uint32
}

const (
	STATE_NORMAL   int32 = 0
	STATE_FOCUSED        = 1
	STATE_PRESSED        = 2
	STATE_DISABLED       = 3
)

type GuiState = int32

const (
	TEXT_ALIGN_LEFT   int32 = 0
	TEXT_ALIGN_CENTER       = 1
	TEXT_ALIGN_RIGHT        = 2
)

type GuiTextAlignment = int32

const (
	DEFAULT     int32 = 0
	LABEL             = 1
	BUTTON            = 2
	TOGGLE            = 3
	SLIDER            = 4
	PROGRESSBAR       = 5
	CHECKBOX          = 6
	COMBOBOX          = 7
	DROPDOWNBOX       = 8
	TEXTBOX           = 9
	VALUEBOX          = 10
	SPINNER           = 11
	LISTVIEW          = 12
	COLORPICKER       = 13
	SCROLLBAR         = 14
	STATUSBAR         = 15
)

type GuiControl = int32

const (
	BORDER_COLOR_NORMAL   int32 = 0
	BASE_COLOR_NORMAL           = 1
	TEXT_COLOR_NORMAL           = 2
	BORDER_COLOR_FOCUSED        = 3
	BASE_COLOR_FOCUSED          = 4
	TEXT_COLOR_FOCUSED          = 5
	BORDER_COLOR_PRESSED        = 6
	BASE_COLOR_PRESSED          = 7
	TEXT_COLOR_PRESSED          = 8
	BORDER_COLOR_DISABLED       = 9
	BASE_COLOR_DISABLED         = 10
	TEXT_COLOR_DISABLED         = 11
	BORDER_WIDTH                = 12
	TEXT_PADDING                = 13
	TEXT_ALIGNMENT              = 14
	RESERVED                    = 15
)

type GuiControlProperty = int32

const (
	TEXT_SIZE        int32 = 16
	TEXT_SPACING           = 17
	LINE_COLOR             = 18
	BACKGROUND_COLOR       = 19
)

type GuiDefaultProperty = int32

const (
	GROUP_PADDING int32 = 16
)

type GuiToggleProperty = int32

const (
	SLIDER_WIDTH   int32 = 16
	SLIDER_PADDING       = 17
)

type GuiSliderProperty = int32

const (
	PROGRESS_PADDING int32 = 16
)

type GuiProgressBarProperty = int32

const (
	ARROWS_SIZE           int32 = 16
	ARROWS_VISIBLE              = 17
	SCROLL_SLIDER_PADDING       = 18
	SCROLL_SLIDER_SIZE          = 19
	SCROLL_PADDING              = 20
	SCROLL_SPEED                = 21
)

type GuiScrollBarProperty = int32

const (
	CHECK_PADDING int32 = 16
)

type GuiCheckBoxProperty = int32

const (
	COMBO_BUTTON_WIDTH   int32 = 16
	COMBO_BUTTON_SPACING       = 17
)

type GuiComboBoxProperty = int32

const (
	ARROW_PADDING          int32 = 16
	DROPDOWN_ITEMS_SPACING       = 17
)

type GuiDropdownBoxProperty = int32

const (
	TEXT_INNER_PADDING int32 = 16
	TEXT_LINES_SPACING       = 17
)

type GuiTextBoxProperty = int32

const (
	SPIN_BUTTON_WIDTH   int32 = 16
	SPIN_BUTTON_SPACING       = 17
)

type GuiSpinnerProperty = int32

const (
	LIST_ITEMS_HEIGHT  int32 = 16
	LIST_ITEMS_SPACING       = 17
	SCROLLBAR_WIDTH          = 18
	SCROLLBAR_SIDE           = 19
)

type GuiListViewProperty = int32

const (
	COLOR_SELECTOR_SIZE      int32 = 16
	HUEBAR_WIDTH                   = 17
	HUEBAR_PADDING                 = 18
	HUEBAR_SELECTOR_HEIGHT         = 19
	HUEBAR_SELECTOR_OVERFLOW       = 20
)

type GuiColorPickerProperty = int32

const (
	ICON_NONE                    int32 = 0
	ICON_FOLDER_FILE_OPEN              = 1
	ICON_FILE_SAVE_CLASSIC             = 2
	ICON_FOLDER_OPEN                   = 3
	ICON_FOLDER_SAVE                   = 4
	ICON_FILE_OPEN                     = 5
	ICON_FILE_SAVE                     = 6
	ICON_FILE_EXPORT                   = 7
	ICON_FILE_ADD                      = 8
	ICON_FILE_DELETE                   = 9
	ICON_FILETYPE_TEXT                 = 10
	ICON_FILETYPE_AUDIO                = 11
	ICON_FILETYPE_IMAGE                = 12
	ICON_FILETYPE_PLAY                 = 13
	ICON_FILETYPE_VIDEO                = 14
	ICON_FILETYPE_INFO                 = 15
	ICON_FILE_COPY                     = 16
	ICON_FILE_CUT                      = 17
	ICON_FILE_PASTE                    = 18
	ICON_CURSOR_HAND                   = 19
	ICON_CURSOR_POINTER                = 20
	ICON_CURSOR_CLASSIC                = 21
	ICON_PENCIL                        = 22
	ICON_PENCIL_BIG                    = 23
	ICON_BRUSH_CLASSIC                 = 24
	ICON_BRUSH_PAINTER                 = 25
	ICON_WATER_DROP                    = 26
	ICON_COLOR_PICKER                  = 27
	ICON_RUBBER                        = 28
	ICON_COLOR_BUCKET                  = 29
	ICON_TEXT_T                        = 30
	ICON_TEXT_A                        = 31
	ICON_SCALE                         = 32
	ICON_RESIZE                        = 33
	ICON_FILTER_POINT                  = 34
	ICON_FILTER_BILINEAR               = 35
	ICON_CROP                          = 36
	ICON_CROP_ALPHA                    = 37
	ICON_SQUARE_TOGGLE                 = 38
	ICON_SYMMETRY                      = 39
	ICON_SYMMETRY_HORIZONTAL           = 40
	ICON_SYMMETRY_VERTICAL             = 41
	ICON_LENS                          = 42
	ICON_LENS_BIG                      = 43
	ICON_EYE_ON                        = 44
	ICON_EYE_OFF                       = 45
	ICON_FILTER_TOP                    = 46
	ICON_FILTER                        = 47
	ICON_TARGET_POINT                  = 48
	ICON_TARGET_SMALL                  = 49
	ICON_TARGET_BIG                    = 50
	ICON_TARGET_MOVE                   = 51
	ICON_CURSOR_MOVE                   = 52
	ICON_CURSOR_SCALE                  = 53
	ICON_CURSOR_SCALE_RIGHT            = 54
	ICON_CURSOR_SCALE_LEFT             = 55
	ICON_UNDO                          = 56
	ICON_REDO                          = 57
	ICON_REREDO                        = 58
	ICON_MUTATE                        = 59
	ICON_ROTATE                        = 60
	ICON_REPEAT                        = 61
	ICON_SHUFFLE                       = 62
	ICON_EMPTYBOX                      = 63
	ICON_TARGET                        = 64
	ICON_TARGET_SMALL_FILL             = 65
	ICON_TARGET_BIG_FILL               = 66
	ICON_TARGET_MOVE_FILL              = 67
	ICON_CURSOR_MOVE_FILL              = 68
	ICON_CURSOR_SCALE_FILL             = 69
	ICON_CURSOR_SCALE_RIGHT_FILL       = 70
	ICON_CURSOR_SCALE_LEFT_FILL        = 71
	ICON_UNDO_FILL                     = 72
	ICON_REDO_FILL                     = 73
	ICON_REREDO_FILL                   = 74
	ICON_MUTATE_FILL                   = 75
	ICON_ROTATE_FILL                   = 76
	ICON_REPEAT_FILL                   = 77
	ICON_SHUFFLE_FILL                  = 78
	ICON_EMPTYBOX_SMALL                = 79
	ICON_BOX                           = 80
	ICON_BOX_TOP                       = 81
	ICON_BOX_TOP_RIGHT                 = 82
	ICON_BOX_RIGHT                     = 83
	ICON_BOX_BOTTOM_RIGHT              = 84
	ICON_BOX_BOTTOM                    = 85
	ICON_BOX_BOTTOM_LEFT               = 86
	ICON_BOX_LEFT                      = 87
	ICON_BOX_TOP_LEFT                  = 88
	ICON_BOX_CENTER                    = 89
	ICON_BOX_CIRCLE_MASK               = 90
	ICON_POT                           = 91
	ICON_ALPHA_MULTIPLY                = 92
	ICON_ALPHA_CLEAR                   = 93
	ICON_DITHERING                     = 94
	ICON_MIPMAPS                       = 95
	ICON_BOX_GRID                      = 96
	ICON_GRID                          = 97
	ICON_BOX_CORNERS_SMALL             = 98
	ICON_BOX_CORNERS_BIG               = 99
	ICON_FOUR_BOXES                    = 100
	ICON_GRID_FILL                     = 101
	ICON_BOX_MULTISIZE                 = 102
	ICON_ZOOM_SMALL                    = 103
	ICON_ZOOM_MEDIUM                   = 104
	ICON_ZOOM_BIG                      = 105
	ICON_ZOOM_ALL                      = 106
	ICON_ZOOM_CENTER                   = 107
	ICON_BOX_DOTS_SMALL                = 108
	ICON_BOX_DOTS_BIG                  = 109
	ICON_BOX_CONCENTRIC                = 110
	ICON_BOX_GRID_BIG                  = 111
	ICON_OK_TICK                       = 112
	ICON_CROSS                         = 113
	ICON_ARROW_LEFT                    = 114
	ICON_ARROW_RIGHT                   = 115
	ICON_ARROW_DOWN                    = 116
	ICON_ARROW_UP                      = 117
	ICON_ARROW_LEFT_FILL               = 118
	ICON_ARROW_RIGHT_FILL              = 119
	ICON_ARROW_DOWN_FILL               = 120
	ICON_ARROW_UP_FILL                 = 121
	ICON_AUDIO                         = 122
	ICON_FX                            = 123
	ICON_WAVE                          = 124
	ICON_WAVE_SINUS                    = 125
	ICON_WAVE_SQUARE                   = 126
	ICON_WAVE_TRIANGULAR               = 127
	ICON_CROSS_SMALL                   = 128
	ICON_PLAYER_PREVIOUS               = 129
	ICON_PLAYER_PLAY_BACK              = 130
	ICON_PLAYER_PLAY                   = 131
	ICON_PLAYER_PAUSE                  = 132
	ICON_PLAYER_STOP                   = 133
	ICON_PLAYER_NEXT                   = 134
	ICON_PLAYER_RECORD                 = 135
	ICON_MAGNET                        = 136
	ICON_LOCK_CLOSE                    = 137
	ICON_LOCK_OPEN                     = 138
	ICON_CLOCK                         = 139
	ICON_TOOLS                         = 140
	ICON_GEAR                          = 141
	ICON_GEAR_BIG                      = 142
	ICON_BIN                           = 143
	ICON_HAND_POINTER                  = 144
	ICON_LASER                         = 145
	ICON_COIN                          = 146
	ICON_EXPLOSION                     = 147
	ICON_1UP                           = 148
	ICON_PLAYER                        = 149
	ICON_PLAYER_JUMP                   = 150
	ICON_KEY                           = 151
	ICON_DEMON                         = 152
	ICON_TEXT_POPUP                    = 153
	ICON_GEAR_EX                       = 154
	ICON_CRACK                         = 155
	ICON_CRACK_POINTS                  = 156
	ICON_STAR                          = 157
	ICON_DOOR                          = 158
	ICON_EXIT                          = 159
	ICON_MODE_2D                       = 160
	ICON_MODE_3D                       = 161
	ICON_CUBE                          = 162
	ICON_CUBE_FACE_TOP                 = 163
	ICON_CUBE_FACE_LEFT                = 164
	ICON_CUBE_FACE_FRONT               = 165
	ICON_CUBE_FACE_BOTTOM              = 166
	ICON_CUBE_FACE_RIGHT               = 167
	ICON_CUBE_FACE_BACK                = 168
	ICON_CAMERA                        = 169
	ICON_SPECIAL                       = 170
	ICON_LINK_NET                      = 171
	ICON_LINK_BOXES                    = 172
	ICON_LINK_MULTI                    = 173
	ICON_LINK                          = 174
	ICON_LINK_BROKE                    = 175
	ICON_TEXT_NOTES                    = 176
	ICON_NOTEBOOK                      = 177
	ICON_SUITCASE                      = 178
	ICON_SUITCASE_ZIP                  = 179
	ICON_MAILBOX                       = 180
	ICON_MONITOR                       = 181
	ICON_PRINTER                       = 182
	ICON_PHOTO_CAMERA                  = 183
	ICON_PHOTO_CAMERA_FLASH            = 184
	ICON_HOUSE                         = 185
	ICON_HEART                         = 186
	ICON_CORNER                        = 187
	ICON_VERTICAL_BARS                 = 188
	ICON_VERTICAL_BARS_FILL            = 189
	ICON_LIFE_BARS                     = 190
	ICON_INFO                          = 191
	ICON_CROSSLINE                     = 192
	ICON_HELP                          = 193
	ICON_FILETYPE_ALPHA                = 194
	ICON_FILETYPE_HOME                 = 195
	ICON_LAYERS_VISIBLE                = 196
	ICON_LAYERS                        = 197
	ICON_WINDOW                        = 198
	ICON_HIDPI                         = 199
	ICON_FILETYPE_BINARY               = 200
	ICON_HEX                           = 201
	ICON_SHIELD                        = 202
	ICON_FILE_NEW                      = 203
	ICON_FOLDER_ADD                    = 204
	ICON_ALARM                         = 205
	ICON_CPU                           = 206
	ICON_ROM                           = 207
	ICON_STEP_OVER                     = 208
	ICON_STEP_INTO                     = 209
	ICON_STEP_OUT                      = 210
	ICON_RESTART                       = 211
	ICON_BREAKPOINT_ON                 = 212
	ICON_BREAKPOINT_OFF                = 213
	ICON_BURGER_MENU                   = 214
	ICON_CASE_SENSITIVE                = 215
	ICON_REG_EXP                       = 216
	ICON_FOLDER                        = 217
	ICON_FILE                          = 218
	ICON_219                           = 219
	ICON_220                           = 220
	ICON_221                           = 221
	ICON_222                           = 222
	ICON_223                           = 223
	ICON_224                           = 224
	ICON_225                           = 225
	ICON_226                           = 226
	ICON_227                           = 227
	ICON_228                           = 228
	ICON_229                           = 229
	ICON_230                           = 230
	ICON_231                           = 231
	ICON_232                           = 232
	ICON_233                           = 233
	ICON_234                           = 234
	ICON_235                           = 235
	ICON_236                           = 236
	ICON_237                           = 237
	ICON_238                           = 238
	ICON_239                           = 239
	ICON_240                           = 240
	ICON_241                           = 241
	ICON_242                           = 242
	ICON_243                           = 243
	ICON_244                           = 244
	ICON_245                           = 245
	ICON_246                           = 246
	ICON_247                           = 247
	ICON_248                           = 248
	ICON_249                           = 249
	ICON_250                           = 250
	ICON_251                           = 251
	ICON_252                           = 252
	ICON_253                           = 253
	ICON_254                           = 254
	ICON_255                           = 255
)

type GuiIconName = int32
type _Bool int32

type va_list struct {
	position int
	Slice    []interface{}
}

func create_va_list(list []interface{}) *va_list {
	return &va_list{
		position: 0,
		Slice:    list,
	}
}

func va_start(v *va_list, count interface{}) {
	v.position = 0
}

func va_end(v *va_list) {
	// do nothing
}

func va_arg(v *va_list) interface{} {
	defer func() {
		v.position++
	}()
	value := v.Slice[v.position]
	switch value.(type) {
	case int:
		return int32(value.(int))
	default:
		return value
	}
}
