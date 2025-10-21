package gleam

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	_ "github.com/kr/pretty"
)

var (
	gleamTestExt = "_test.gleam"
	gleamExt     = ".gleam"
	erlExt       = ".erl"

	errSkipImport    errorType = "skip"
	errNotFound      errorType = "not found"
	errMultipleFound errorType = "multiple found"

	// https://www.erlang.org/doc/man_index.html
	// Global Erlang interop modules
	erlangStdlibModules = map[string]bool{
		"alarm_handler":                  true,
		"application":                    true,
		"argparse":                       true,
		"array":                          true,
		"asn1ct":                         true,
		"atomics":                        true,
		"auth":                           true,
		"base64":                         true,
		"beam_lib":                       true,
		"binary":                         true,
		"c":                              true,
		"calendar":                       true,
		"cerl":                           true,
		"cerl_clauses":                   true,
		"cerl_trees":                     true,
		"code":                           true,
		"compile":                        true,
		"counters":                       true,
		"cover":                          true,
		"cprof":                          true,
		"cpu_sup":                        true,
		"crashdump_viewer":               true,
		"crypto":                         true,
		"ct":                             true,
		"ct_cover":                       true,
		"ct_ftp":                         true,
		"ct_hooks":                       true,
		"ct_master":                      true,
		"ct_netconfc":                    true,
		"ct_property_test":               true,
		"ct_rpc":                         true,
		"ct_slave":                       true,
		"ct_snmp":                        true,
		"ct_ssh":                         true,
		"ct_suite":                       true,
		"ct_telnet":                      true,
		"ct_testspec":                    true,
		"dbg":                            true,
		"debugger":                       true,
		"dets":                           true,
		"dialyzer":                       true,
		"diameter":                       true,
		"diameter_app":                   true,
		"diameter_codec":                 true,
		"diameter_make":                  true,
		"diameter_sctp":                  true,
		"diameter_tcp":                   true,
		"diameter_transport":             true,
		"dict":                           true,
		"digraph":                        true,
		"digraph_utils":                  true,
		"disk_log":                       true,
		"disksup":                        true,
		"dyntrace":                       true,
		"edlin":                          true,
		"edlin_expand":                   true,
		"edoc":                           true,
		"edoc_doclet":                    true,
		"edoc_doclet_chunks":             true,
		"edoc_doclet_markdown":           true,
		"edoc_extract":                   true,
		"edoc_html_to_markdown":          true,
		"edoc_layout":                    true,
		"edoc_layout_chunks":             true,
		"edoc_lib":                       true,
		"edoc_run":                       true,
		"eldap":                          true,
		"epp":                            true,
		"epp_dodger":                     true,
		"eprof":                          true,
		"erl_anno":                       true,
		"erl_boot_server":                true,
		"erl_comment_scan":               true,
		"erl_ddll":                       true,
		"erl_debugger":                   true,
		"erl_epmd":                       true,
		"erl_error":                      true,
		"erl_eval":                       true,
		"erl_expand_records":             true,
		"erl_features":                   true,
		"erl_id_trans":                   true,
		"erl_internal":                   true,
		"erl_lint":                       true,
		"erl_parse":                      true,
		"erl_pp":                         true,
		"erl_prettypr":                   true,
		"erl_prim_loader":                true,
		"erl_recomment":                  true,
		"erl_scan":                       true,
		"erl_syntax":                     true,
		"erl_syntax_lib":                 true,
		"erl_tar":                        true,
		"erl_tracer":                     true,
		"erlang":                         true,
		"erpc":                           true,
		"error_handler":                  true,
		"error_logger":                   true,
		"escript":                        true,
		"et":                             true,
		"et_collector":                   true,
		"et_selector":                    true,
		"et_viewer":                      true,
		"etop":                           true,
		"ets":                            true,
		"eunit":                          true,
		"eunit_surefire":                 true,
		"file":                           true,
		"file_sorter":                    true,
		"filelib":                        true,
		"filename":                       true,
		"fprof":                          true,
		"ftp":                            true,
		"gb_sets":                        true,
		"gb_trees":                       true,
		"gen_event":                      true,
		"gen_fsm":                        true,
		"gen_sctp":                       true,
		"gen_server":                     true,
		"gen_statem":                     true,
		"gen_tcp":                        true,
		"gen_udp":                        true,
		"gl":                             true,
		"global":                         true,
		"global_group":                   true,
		"glu":                            true,
		"heart":                          true,
		"http_uri":                       true,
		"httpc":                          true,
		"httpd":                          true,
		"httpd_custom_api":               true,
		"httpd_socket":                   true,
		"httpd_util":                     true,
		"i":                              true,
		"inet":                           true,
		"inet_res":                       true,
		"inets":                          true,
		"init":                           true,
		"instrument":                     true,
		"int":                            true,
		"io":                             true,
		"io_lib":                         true,
		"json":                           true,
		"lcnt":                           true,
		"leex":                           true,
		"lists":                          true,
		"log_mf_h":                       true,
		"logger":                         true,
		"logger_disk_log_h":              true,
		"logger_filters":                 true,
		"logger_formatter":               true,
		"logger_handler":                 true,
		"logger_std_h":                   true,
		"make":                           true,
		"maps":                           true,
		"math":                           true,
		"megaco":                         true,
		"megaco_codec_meas":              true,
		"megaco_codec_mstone1":           true,
		"megaco_codec_mstone2":           true,
		"megaco_codec_transform":         true,
		"megaco_digit_map":               true,
		"megaco_edist_compress":          true,
		"megaco_encoder":                 true,
		"megaco_flex_scanner":            true,
		"megaco_sdp":                     true,
		"megaco_tcp":                     true,
		"megaco_transport":               true,
		"megaco_udp":                     true,
		"megaco_user":                    true,
		"memsup":                         true,
		"merl":                           true,
		"merl_transform":                 true,
		"mnesia":                         true,
		"mnesia_frag_hash":               true,
		"mnesia_registry":                true,
		"mod_alias":                      true,
		"mod_auth":                       true,
		"mod_esi":                        true,
		"mod_security":                   true,
		"ms_transform":                   true,
		"msacc":                          true,
		"net":                            true,
		"net_adm":                        true,
		"net_kernel":                     true,
		"nteventlog":                     true,
		"observer":                       true,
		"odbc":                           true,
		"orddict":                        true,
		"ordsets":                        true,
		"os":                             true,
		"os_sup":                         true,
		"peer":                           true,
		"persistent_term":                true,
		"pg":                             true,
		"pool":                           true,
		"prettypr":                       true,
		"proc_lib":                       true,
		"proplists":                      true,
		"public_key":                     true,
		"qlc":                            true,
		"queue":                          true,
		"rand":                           true,
		"random":                         true,
		"rb":                             true,
		"re":                             true,
		"release_handler":                true,
		"reltool":                        true,
		"rpc":                            true,
		"scheduler":                      true,
		"seq_trace":                      true,
		"sets":                           true,
		"shell":                          true,
		"shell_default":                  true,
		"shell_docs":                     true,
		"slave":                          true,
		"snmp":                           true,
		"snmp_community_mib":             true,
		"snmp_framework_mib":             true,
		"snmp_generic":                   true,
		"snmp_index":                     true,
		"snmp_notification_mib":          true,
		"snmp_pdus":                      true,
		"snmp_standard_mib":              true,
		"snmp_target_mib":                true,
		"snmp_user_based_sm_mib":         true,
		"snmp_view_based_acm_mib":        true,
		"snmpa":                          true,
		"snmpa_conf":                     true,
		"snmpa_discovery_handler":        true,
		"snmpa_error":                    true,
		"snmpa_error_io":                 true,
		"snmpa_error_logger":             true,
		"snmpa_error_report":             true,
		"snmpa_local_db":                 true,
		"snmpa_mib_data":                 true,
		"snmpa_mib_storage":              true,
		"snmpa_mpd":                      true,
		"snmpa_network_interface":        true,
		"snmpa_network_interface_filter": true,
		"snmpa_notification_delivery_info_receiver": true,
		"snmpa_notification_filter":                 true,
		"snmpa_supervisor":                          true,
		"snmpc":                                     true,
		"snmpm":                                     true,
		"snmpm_conf":                                true,
		"snmpm_mpd":                                 true,
		"snmpm_network_interface":                   true,
		"snmpm_network_interface_filter":            true,
		"snmpm_user":                                true,
		"socket":                                    true,
		"sofs":                                      true,
		"ssh":                                       true,
		"ssh_agent":                                 true,
		"ssh_client_channel":                        true,
		"ssh_client_key_api":                        true,
		"ssh_connection":                            true,
		"ssh_file":                                  true,
		"ssh_server_channel":                        true,
		"ssh_server_key_api":                        true,
		"ssh_sftp":                                  true,
		"ssh_sftpd":                                 true,
		"ssl":                                       true,
		"ssl_crl_cache":                             true,
		"ssl_crl_cache_api":                         true,
		"ssl_session_cache_api":                     true,
		"string":                                    true,
		"supervisor":                                true,
		"supervisor_bridge":                         true,
		"sys":                                       true,
		"system_information":                        true,
		"systools":                                  true,
		"tags":                                      true,
		"tcp":                                       true,
		"tftp":                                      true,
		"tftp_logger":                               true,
		"timer":                                     true,
		"tprof":                                     true,
		"trace":                                     true,
		"ttb":                                       true,
		"unicode":                                   true,
		"unix_telnet":                               true,
		"uri_string":                                true,
		"win32reg":                                  true,
		"wrap_log_reader":                           true,
		"wx":                                        true,
		"wxAcceleratorEntry":                        true,
		"wxAcceleratorTable":                        true,
		"wxActivateEvent":                           true,
		"wxArtProvider":                             true,
		"wxAuiDockArt":                              true,
		"wxAuiManager":                              true,
		"wxAuiManagerEvent":                         true,
		"wxAuiNotebook":                             true,
		"wxAuiNotebookEvent":                        true,
		"wxAuiPaneInfo":                             true,
		"wxAuiSimpleTabArt":                         true,
		"wxAuiTabArt":                               true,
		"wxBitmap":                                  true,
		"wxBitmapButton":                            true,
		"wxBitmapDataObject":                        true,
		"wxBookCtrlBase":                            true,
		"wxBookCtrlEvent":                           true,
		"wxBoxSizer":                                true,
		"wxBrush":                                   true,
		"wxBufferedDC":                              true,
		"wxBufferedPaintDC":                         true,
		"wxButton":                                  true,
		"wxCalendarCtrl":                            true,
		"wxCalendarDateAttr":                        true,
		"wxCalendarEvent":                           true,
		"wxCaret":                                   true,
		"wxCheckBox":                                true,
		"wxCheckListBox":                            true,
		"wxChildFocusEvent":                         true,
		"wxChoice":                                  true,
		"wxChoicebook":                              true,
		"wxClientDC":                                true,
		"wxClipboard":                               true,
		"wxClipboardTextEvent":                      true,
		"wxCloseEvent":                              true,
		"wxColourData":                              true,
		"wxColourDialog":                            true,
		"wxColourPickerCtrl":                        true,
		"wxColourPickerEvent":                       true,
		"wxComboBox":                                true,
		"wxCommandEvent":                            true,
		"wxContextMenuEvent":                        true,
		"wxControl":                                 true,
		"wxControlWithItems":                        true,
		"wxCursor":                                  true,
		"wxDC":                                      true,
		"wxDCOverlay":                               true,
		"wxDataObject":                              true,
		"wxDateEvent":                               true,
		"wxDatePickerCtrl":                          true,
		"wxDialog":                                  true,
		"wxDirDialog":                               true,
		"wxDirPickerCtrl":                           true,
		"wxDisplay":                                 true,
		"wxDisplayChangedEvent":                     true,
		"wxDropFilesEvent":                          true,
		"wxEraseEvent":                              true,
		"wxEvent":                                   true,
		"wxEvtHandler":                              true,
		"wxFileDataObject":                          true,
		"wxFileDialog":                              true,
		"wxFileDirPickerEvent":                      true,
		"wxFilePickerCtrl":                          true,
		"wxFindReplaceData":                         true,
		"wxFindReplaceDialog":                       true,
		"wxFlexGridSizer":                           true,
		"wxFocusEvent":                              true,
		"wxFont":                                    true,
		"wxFontData":                                true,
		"wxFontDialog":                              true,
		"wxFontPickerCtrl":                          true,
		"wxFontPickerEvent":                         true,
		"wxFrame":                                   true,
		"wxGBSizerItem":                             true,
		"wxGCDC":                                    true,
		"wxGLCanvas":                                true,
		"wxGLContext":                               true,
		"wxGauge":                                   true,
		"wxGenericDirCtrl":                          true,
		"wxGraphicsBrush":                           true,
		"wxGraphicsContext":                         true,
		"wxGraphicsFont":                            true,
		"wxGraphicsGradientStops":                   true,
		"wxGraphicsMatrix":                          true,
		"wxGraphicsObject":                          true,
		"wxGraphicsPath":                            true,
		"wxGraphicsPen":                             true,
		"wxGraphicsRenderer":                        true,
		"wxGrid":                                    true,
		"wxGridBagSizer":                            true,
		"wxGridCellAttr":                            true,
		"wxGridCellBoolEditor":                      true,
		"wxGridCellBoolRenderer":                    true,
		"wxGridCellChoiceEditor":                    true,
		"wxGridCellEditor":                          true,
		"wxGridCellFloatEditor":                     true,
		"wxGridCellFloatRenderer":                   true,
		"wxGridCellNumberEditor":                    true,
		"wxGridCellNumberRenderer":                  true,
		"wxGridCellRenderer":                        true,
		"wxGridCellStringRenderer":                  true,
		"wxGridCellTextEditor":                      true,
		"wxGridEvent":                               true,
		"wxGridSizer":                               true,
		"wxHelpEvent":                               true,
		"wxHtmlEasyPrinting":                        true,
		"wxHtmlLinkEvent":                           true,
		"wxHtmlWindow":                              true,
		"wxIcon":                                    true,
		"wxIconBundle":                              true,
		"wxIconizeEvent":                            true,
		"wxIdleEvent":                               true,
		"wxImage":                                   true,
		"wxImageList":                               true,
		"wxInitDialogEvent":                         true,
		"wxJoystickEvent":                           true,
		"wxKeyEvent":                                true,
		"wxLayoutAlgorithm":                         true,
		"wxListBox":                                 true,
		"wxListCtrl":                                true,
		"wxListEvent":                               true,
		"wxListItem":                                true,
		"wxListItemAttr":                            true,
		"wxListView":                                true,
		"wxListbook":                                true,
		"wxLocale":                                  true,
		"wxLogNull":                                 true,
		"wxMDIChildFrame":                           true,
		"wxMDIClientWindow":                         true,
		"wxMDIParentFrame":                          true,
		"wxMask":                                    true,
		"wxMaximizeEvent":                           true,
		"wxMemoryDC":                                true,
		"wxMenu":                                    true,
		"wxMenuBar":                                 true,
		"wxMenuEvent":                               true,
		"wxMenuItem":                                true,
		"wxMessageDialog":                           true,
		"wxMiniFrame":                               true,
		"wxMirrorDC":                                true,
		"wxMouseCaptureChangedEvent":                true,
		"wxMouseCaptureLostEvent":                   true,
		"wxMouseEvent":                              true,
		"wxMoveEvent":                               true,
		"wxMultiChoiceDialog":                       true,
		"wxNavigationKeyEvent":                      true,
		"wxNotebook":                                true,
		"wxNotificationMessage":                     true,
		"wxNotifyEvent":                             true,
		"wxOverlay":                                 true,
		"wxPageSetupDialog":                         true,
		"wxPageSetupDialogData":                     true,
		"wxPaintDC":                                 true,
		"wxPaintEvent":                              true,
		"wxPalette":                                 true,
		"wxPaletteChangedEvent":                     true,
		"wxPanel":                                   true,
		"wxPasswordEntryDialog":                     true,
		"wxPen":                                     true,
		"wxPickerBase":                              true,
		"wxPopupTransientWindow":                    true,
		"wxPopupWindow":                             true,
		"wxPostScriptDC":                            true,
		"wxPreviewCanvas":                           true,
		"wxPreviewControlBar":                       true,
		"wxPreviewFrame":                            true,
		"wxPrintData":                               true,
		"wxPrintDialog":                             true,
		"wxPrintDialogData":                         true,
		"wxPrintPreview":                            true,
		"wxPrinter":                                 true,
		"wxPrintout":                                true,
		"wxProgressDialog":                          true,
		"wxQueryNewPaletteEvent":                    true,
		"wxRadioBox":                                true,
		"wxRadioButton":                             true,
		"wxRegion":                                  true,
		"wxSashEvent":                               true,
		"wxSashLayoutWindow":                        true,
		"wxSashWindow":                              true,
		"wxScreenDC":                                true,
		"wxScrollBar":                               true,
		"wxScrollEvent":                             true,
		"wxScrollWinEvent":                          true,
		"wxScrolledWindow":                          true,
		"wxSetCursorEvent":                          true,
		"wxShowEvent":                               true,
		"wxSingleChoiceDialog":                      true,
		"wxSizeEvent":                               true,
		"wxSizer":                                   true,
		"wxSizerFlags":                              true,
		"wxSizerItem":                               true,
		"wxSlider":                                  true,
		"wxSpinButton":                              true,
		"wxSpinCtrl":                                true,
		"wxSpinEvent":                               true,
		"wxSplashScreen":                            true,
		"wxSplitterEvent":                           true,
		"wxSplitterWindow":                          true,
		"wxStaticBitmap":                            true,
		"wxStaticBox":                               true,
		"wxStaticBoxSizer":                          true,
		"wxStaticLine":                              true,
		"wxStaticText":                              true,
		"wxStatusBar":                               true,
		"wxStdDialogButtonSizer":                    true,
		"wxStyledTextCtrl":                          true,
		"wxStyledTextEvent":                         true,
		"wxSysColourChangedEvent":                   true,
		"wxSystemOptions":                           true,
		"wxSystemSettings":                          true,
		"wxTaskBarIcon":                             true,
		"wxTaskBarIconEvent":                        true,
		"wxTextAttr":                                true,
		"wxTextCtrl":                                true,
		"wxTextDataObject":                          true,
		"wxTextEntryDialog":                         true,
		"wxToggleButton":                            true,
		"wxToolBar":                                 true,
		"wxToolTip":                                 true,
		"wxToolbook":                                true,
		"wxTopLevelWindow":                          true,
		"wxTreeCtrl":                                true,
		"wxTreeEvent":                               true,
		"wxTreebook":                                true,
		"wxUpdateUIEvent":                           true,
		"wxWebView":                                 true,
		"wxWebViewEvent":                            true,
		"wxWindow":                                  true,
		"wxWindowCreateEvent":                       true,
		"wxWindowDC":                                true,
		"wxWindowDestroyEvent":                      true,
		"wxXmlResource":                             true,
		"wx_misc":                                   true,
		"wx_object":                                 true,
		"xmerl":                                     true,
		"xmerl_eventp":                              true,
		"xmerl_sax_parser":                          true,
		"xmerl_scan":                                true,
		"xmerl_xpath":                               true,
		"xmerl_xs":                                  true,
		"xmerl_xsd":                                 true,
		"xref":                                      true,
		"yecc":                                      true,
		"zip":                                       true,
		"zlib":                                      true,
		"zstd":                                      true,
	}
)

type errorType string

type gleamGazelleError struct {
	msg       string
	errorType errorType
}

func (gge *gleamGazelleError) ErrorType() errorType {
	return gge.errorType
}

func (gge *gleamGazelleError) Error() string {
	return gge.msg
}

func (*gleamLanguage) Name() string { return "gleam" }

// Returns the Gleam specific import path for the rule.
// These are all of the Gleam modules, declared in srcs.
func (g *gleamLanguage) Imports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
	if !isGleamLibrary(r) {
		return nil
	}

	imports := []resolve.ImportSpec{}
	for _, src := range r.AttrStrings("srcs") {
		if path.Ext(src) == gleamExt {
			imports = append(imports, resolve.ImportSpec{Lang: g.Name(), Imp: path.Join(f.Pkg, strings.TrimSuffix(src, gleamExt))})
		} else if path.Ext(src) == erlExt {
			imports = append(imports, resolve.ImportSpec{Lang: g.Name(), Imp: strings.Join([]string{"erl", strings.TrimSuffix(src, erlExt)}, ":")})
		}
	}

	return imports
}

func (g *gleamLanguage) Embeds(r *rule.Rule, from label.Label) []label.Label {
	return nil
}

// Resolve adds deps to the given rule based on the importRaws.
func (g *gleamLanguage) Resolve(c *config.Config, ix *resolve.RuleIndex, _rc *repo.RemoteCache, r *rule.Rule, importRaws interface{}, from label.Label) {
	// If no imports are given, bail early.
	if importRaws == nil {
		return
	}

	var err error
	gleamConfig := GetGleamConfig(c)
	rc, cleanup := repo.NewRemoteCache(gleamConfig.repos)
	defer func() {
		if cerr := cleanup(); err == nil && cerr != nil {
			err = cerr
		}
	}()

	imports := importRaws.([]string)
	r.DelAttr("deps")

	// Create a set of dependencies so we can avoid duplicates.
	depSet := make(map[string]bool)
	for _, imp := range imports {
		depLabel, err := g.resolveGleam(c, ix, rc, r, imp, from)
		if err != nil && err.ErrorType() == errSkipImport {
			// If resolveGleam returns errSkipImport, skip this import.
			continue
		} else if err != nil {
			// If resolveGleam has any other error, log it.
			log.Print(err.msg)
			if gleamConfig.externalRepo {
				panic(fmt.Sprintf("failed to resolve dependency for external package: this should not happen, %s", err.msg))
			}
		} else {
			var label label.Label
			if depLabel.Pkg == from.Pkg && depLabel.Repo == from.Repo {
				label = depLabel.Rel(depLabel.Repo, depLabel.Pkg)
			} else {
				label = depLabel.Abs(depLabel.Repo, depLabel.Pkg)
			}
			depSet[label.String()] = true
		}
	}

	if len(depSet) != 0 {
		// If there are dependencies, set the deps attribute.
		deps := make([]string, 0, len(depSet))
		for dep := range depSet {
			deps = append(deps, dep)
		}
		sort.Strings(deps)
		r.SetAttr("deps", deps)
	}
}

// For gleamlibrary rule that does self import modules in srcs, we don't need labels for these.
func isSelfImport(r *rule.Rule, f label.Label, imp string) bool {
	localImports := asSet(mapper(r.AttrStrings("srcs"), func(src string) string {
		return path.Join(f.Pkg, strings.TrimSuffix(src, gleamExt))
	}))
	return localImports[imp]
}

func (g *gleamLanguage) resolveGleam(c *config.Config, ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, imp string, from label.Label) (label.Label, *gleamGazelleError) {
	if erlangStdlibModules[strings.TrimPrefix(imp, "erl:")] {
		return label.NoLabel, &gleamGazelleError{msg: fmt.Sprintf("erlang stdlib module: %s", imp), errorType: errSkipImport}
	}
	if isSelfImport(r, from, imp) {
		return label.NoLabel, &gleamGazelleError{msg: fmt.Sprintf("self import: %s", imp), errorType: errSkipImport}
	}
	results := ix.FindRulesByImportWithConfig(c, resolve.ImportSpec{Lang: g.Name(), Imp: imp}, g.Name())
	if len(results) == 0 {
		l, err := g.tryResolveExternalDeps(c, ix, rc, r, imp, from)
		if err != nil {
			return label.NoLabel, &gleamGazelleError{msg: fmt.Sprintf("no rule may be imported with %q from package %s: %v", imp, from, err), errorType: errNotFound}
		}
		return l, nil
	} else if len(results) > 1 {
		return label.NoLabel, &gleamGazelleError{msg: fmt.Sprintf("multiple rules (%s and %s) may be imported with %q from %s", results[0].Label, results[1].Label, imp, from), errorType: errMultipleFound}
	}

	return results[0].Label, nil
}

func (g *gleamLanguage) tryResolveExternalDeps(
	c *config.Config,
	ix *resolve.RuleIndex,
	rc *repo.RemoteCache,
	r *rule.Rule,
	imp string,
	from label.Label,
) (label.Label, error) {
	pkg, module, err := rc.Root(imp)
	if err != nil {
		return label.NoLabel, err
	}
	depPkg := filepath.Dir(pkg)
	gleamModule := filepath.Base(pkg)
	if depPkg == "." {
		depPkg = ""
	}
	depMod := externalGetEquivalentGleamImport(gleamModule)
	if strings.Contains(depMod, "/") {
		depPkg = depPkg + filepath.Dir(depMod)
		depMod = filepath.Base(depMod)
	}
	l := label.New(module, depPkg, depMod)
	return l, nil
}

func externalGetEquivalentGleamImport(
	imp string,
) string {
	if strings.HasPrefix(imp, "erl:") {
		erlPath, _ := strings.CutPrefix(imp, "erl:")
		if strings.Contains(imp, "@") {
			return strings.ReplaceAll(erlPath, "@", "/")
		} else {
			return fmt.Sprintf("%s_ffi", erlPath)
		}
	}
	return imp
}
