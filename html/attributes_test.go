package html_test

import (
	"fmt"
	"testing"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"maragu.dev/gomponents/internal/assert"
)

func TestBooleanAttributes(t *testing.T) {
	tests := []struct {
		Name string
		Func func() g.Node
	}{
		{Name: "allowfullscreen", Func: AllowFullscreen},
		{Name: "async", Func: Async},
		{Name: "autofocus", Func: AutoFocus},
		{Name: "autoplay", Func: AutoPlay},
		{Name: "checked", Func: Checked},
		{Name: "controls", Func: Controls},
		{Name: "default", Func: Default},
		{Name: "defer", Func: Defer},
		{Name: "disabled", Func: Disabled},
		{Name: "formnovalidate", Func: FormNoValidate},
		{Name: "inert", Func: Inert},
		{Name: "ismap", Func: IsMap},
		{Name: "itemscope", Func: ItemScope},
		{Name: "loop", Func: Loop},
		{Name: "multiple", Func: Multiple},
		{Name: "muted", Func: Muted},
		{Name: "nomodule", Func: NoModule},
		{Name: "novalidate", Func: NoValidate},
		{Name: "open", Func: Open},
		{Name: "playsinline", Func: PlaysInline},
		{Name: "readonly", Func: ReadOnly},
		{Name: "required", Func: Required},
		{Name: "reversed", Func: Reversed},
		{Name: "selected", Func: Selected},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			n := g.El("div", test.Func())
			assert.Equal(t, fmt.Sprintf(`<div %v></div>`, test.Name), n)
		})
	}
}

func TestSimpleAttributes(t *testing.T) {
	tests := []struct {
		Name string
		Func func(string) g.Node
	}{
		{Name: "abbr", Func: AbbrAttr},
		{Name: "accept", Func: Accept},
		{Name: "accept-charset", Func: AcceptCharset},
		{Name: "accesskey", Func: AccessKey},
		{Name: "action", Func: Action},
		{Name: "align", Func: Align},
		{Name: "allow", Func: Allow},
		{Name: "allowpaymentrequest", Func: AllowPaymentRequest},
		{Name: "alpha", Func: Alpha},
		{Name: "alt", Func: Alt},
		{Name: "archive", Func: Archive},
		{Name: "as", Func: As},
		{Name: "attributionsourceid", Func: AttributionSourceID},
		{Name: "attributionsrc", Func: AttributionSrc},
		{Name: "autocapitalize", Func: AutoCapitalize},
		{Name: "autocomplete", Func: AutoComplete},
		{Name: "autocorrect", Func: AutoCorrect},
		{Name: "axis", Func: Axis},
		{Name: "behavior", Func: Behavior},
		{Name: "blocking", Func: Blocking},
		{Name: "browsingtopics", Func: BrowsingTopics},
		{Name: "capture", Func: Capture},
		{Name: "cellpadding", Func: CellPadding},
		{Name: "cellspacing", Func: CellSpacing},
		{Name: "char", Func: Char},
		{Name: "charoff", Func: CharOff},
		{Name: "charset", Func: Charset},
		{Name: "cite", Func: CiteAttr},
		{Name: "class", Func: Class},
		{Name: "classid", Func: Classid},
		{Name: "clear", Func: Clear},
		{Name: "closedby", Func: Closedby},
		{Name: "codebase", Func: Codebase},
		{Name: "codetype", Func: Codetype},
		{Name: "colorspace", Func: Colorspace},
		{Name: "cols", Func: Cols},
		{Name: "colspan", Func: ColSpan},
		{Name: "command", Func: Command},
		{Name: "commandfor", Func: Commandfor},
		{Name: "compact", Func: Compact},
		{Name: "content", Func: Content},
		{Name: "contenteditable", Func: ContentEditable},
		{Name: "controlslist", Func: ControlsList},
		{Name: "coords", Func: Coords},
		{Name: "credentialless", Func: Credentialless},
		{Name: "crossorigin", Func: CrossOrigin},
		{Name: "csp", Func: Csp},
		{Name: "datetime", Func: DateTime},
		{Name: "declare", Func: Declare},
		{Name: "decoding", Func: Decoding},
		{Name: "dir", Func: Dir},
		{Name: "direction", Func: Direction},
		{Name: "dirname", Func: Dirname},
		{Name: "disablepictureinpicture", Func: DisablePictureInPicture},
		{Name: "disableremoteplayback", Func: DisableRemotePlayback},
		{Name: "download", Func: Download},
		{Name: "draggable", Func: Draggable},
		{Name: "enctype", Func: EncType},
		{Name: "enterkeyhint", Func: EnterKeyHint},
		{Name: "face", Func: Face},
		{Name: "fetchpriority", Func: FetchPriority},
		{Name: "for", Func: For},
		{Name: "form", Func: FormAttr},
		{Name: "formaction", Func: FormAction},
		{Name: "formenctype", Func: FormEncType},
		{Name: "formmethod", Func: FormMethod},
		{Name: "formtarget", Func: FormTarget},
		{Name: "frame", Func: Frame},
		{Name: "frameborder", Func: Frameborder},
		{Name: "headers", Func: Headers},
		{Name: "height", Func: Height},
		{Name: "hidden", Func: Hidden},
		{Name: "high", Func: High},
		{Name: "href", Func: Href},
		{Name: "hreflang", Func: HrefLang},
		{Name: "hreftranslate", Func: HrefTranslate},
		{Name: "hspace", Func: Hspace},
		{Name: "http-equiv", Func: HTTPEquiv},
		{Name: "id", Func: ID},
		{Name: "imagesizes", Func: ImageSizes},
		{Name: "imagesrcset", Func: ImageSrcSet},
		{Name: "inputmode", Func: InputMode},
		{Name: "integrity", Func: Integrity},
		{Name: "interestfor", Func: Interestfor},
		{Name: "is", Func: Is},
		{Name: "itemid", Func: ItemID},
		{Name: "itemprop", Func: ItemProp},
		{Name: "itemref", Func: ItemRef},
		{Name: "itemtype", Func: ItemType},
		{Name: "kind", Func: Kind},
		{Name: "label", Func: LabelAttr},
		{Name: "lang", Func: Lang},
		{Name: "list", Func: List},
		{Name: "loading", Func: Loading},
		{Name: "longdesc", Func: Longdesc},
		{Name: "low", Func: Low},
		{Name: "max", Func: Max},
		{Name: "maxlength", Func: MaxLength},
		{Name: "media", Func: Media},
		{Name: "method", Func: Method},
		{Name: "min", Func: Min},
		{Name: "minlength", Func: MinLength},
		{Name: "name", Func: Name},
		{Name: "nohref", Func: Nohref},
		{Name: "nonce", Func: Nonce},
		{Name: "noresize", Func: Noresize},
		{Name: "noshade", Func: Noshade},
		{Name: "optimum", Func: Optimum},
		{Name: "part", Func: Part},
		{Name: "pattern", Func: Pattern},
		{Name: "ping", Func: Ping},
		{Name: "placeholder", Func: Placeholder},
		{Name: "popovertarget", Func: PopoverTarget},
		{Name: "popovertargetaction", Func: PopoverTargetAction},
		{Name: "poster", Func: Poster},
		{Name: "preload", Func: Preload},
		{Name: "privateToken", Func: PrivateToken},
		{Name: "referrerpolicy", Func: ReferrerPolicy},
		{Name: "rel", Func: Rel},
		{Name: "rev", Func: Rev},
		{Name: "role", Func: Role},
		{Name: "rows", Func: Rows},
		{Name: "rowspan", Func: RowSpan},
		{Name: "rules", Func: Rules},
		{Name: "sandbox", Func: Sandbox},
		{Name: "scheme", Func: Scheme},
		{Name: "scope", Func: Scope},
		{Name: "scrollamount", Func: Scrollamount},
		{Name: "scrolldelay", Func: Scrolldelay},
		{Name: "scrolling", Func: Scrolling},
		{Name: "shadowrootclonable", Func: ShadowRootClonable},
		{Name: "shadowrootdelegatesfocus", Func: ShadowRootDelegatesFocus},
		{Name: "shadowrootmode", Func: ShadowRootMode},
		{Name: "shadowrootserializable", Func: ShadowRootSerializable},
		{Name: "shape", Func: Shape},
		{Name: "size", Func: Size},
		{Name: "sizes", Func: Sizes},
		{Name: "slot", Func: SlotAttr},
		{Name: "span", Func: SpanAttr},
		{Name: "spellcheck", Func: SpellCheck},
		{Name: "src", Func: Src},
		{Name: "srcdoc", Func: SrcDoc},
		{Name: "srclang", Func: SrcLang},
		{Name: "srcset", Func: SrcSet},
		{Name: "standby", Func: Standby},
		{Name: "start", Func: Start},
		{Name: "step", Func: Step},
		{Name: "style", Func: Style},
		{Name: "style", Func: StyleAttr},
		{Name: "summary", Func: SummaryAttr},
		{Name: "tabindex", Func: TabIndex},
		{Name: "target", Func: Target},
		{Name: "title", Func: Title},
		{Name: "title", Func: TitleAttr},
		{Name: "translate", Func: Translate},
		{Name: "truespeed", Func: Truespeed},
		{Name: "type", Func: Type},
		{Name: "usemap", Func: UseMap},
		{Name: "valign", Func: Valign},
		{Name: "value", Func: Value},
		{Name: "valuetype", Func: Valuetype},
		{Name: "version", Func: Version},
		{Name: "virtualkeyboardpolicy", Func: VirtualKeyboardPolicy},
		{Name: "vspace", Func: Vspace},
		{Name: "webkitdirectory", Func: WebkitDirectory},
		{Name: "width", Func: Width},
		{Name: "wrap", Func: Wrap},
		{Name: "writingsuggestions", Func: WritingSuggestions},
		{Name: "xmlns", Func: Xmlns},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			n := g.El("div", test.Func("hat"))
			assert.Equal(t, fmt.Sprintf(`<div %v="hat"></div>`, test.Name), n)
		})
	}
}

func TestVariadicAttributes(t *testing.T) {
	tests := []struct {
		Name string
		Func func(...string) g.Node
	}{
		{Name: "popover", Func: Popover},
	}

	for _, test := range tests {
		t.Run(test.Name + "(no args)", func(t *testing.T) {
			n := g.El("div", test.Func())
			assert.Equal(t, fmt.Sprintf(`<div %v></div>`, test.Name), n)
		})

		t.Run(test.Name +"(one arg)", func(t *testing.T) {
			n := g.El("div", test.Func("hat"))
			assert.Equal(t, fmt.Sprintf(`<div %v="hat"></div>`, test.Name), n)
		})

		t.Run(test.Name + "(two args panics)", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected a panic")
				}
			}()
			n := g.El("div", test.Func("hat", "party"))
			assert.Equal(t, "unreachable", n)
		})
	}
}

func TestAria(t *testing.T) {
	t.Run("returns an attribute which name is prefixed with aria-", func(t *testing.T) {
		n := Aria("selected", "true")
		assert.Equal(t, ` aria-selected="true"`, n)
	})
}

func TestData(t *testing.T) {
	t.Run("returns an attribute which name is prefixed with data-", func(t *testing.T) {
		n := Data("id", "partyhat")
		assert.Equal(t, ` data-id="partyhat"`, n)

		n = DataAttr("id", "partyhat")
		assert.Equal(t, ` data-id="partyhat"`, n)
	})
}
