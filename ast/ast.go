// Package ast parses the clang AST output into AST structures.
package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Konstantin8105/c4go/util"
)

// Node represents any node in the AST.
type Node interface {
	Address() Address
	Children() []Node
	AddChild(node Node)
	Position() Position
}

// Address contains the memory address (originally outputted as a hexadecimal
// string) from the clang AST. The address are not predictable between run and
// are only useful for identifying nodes in a single AST.
//
// The Address is used like a primary key when storing the tree as a flat
// structure.
type Address uint64

// ParseAddress returns the integer representation of the hexadecimal address
// (like 0x7f8a1d8ccfd0). If the address cannot be parsed, 0 is returned.
func ParseAddress(address string) Address {
	addr, _ := strconv.ParseUint(address, 0, 64)

	return Address(addr)
}

// Parse takes the coloured output of the clang AST command and returns a root
// node for the AST.
func Parse(fullline string) (returnNode Node, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Cannot parse line: `%v`. %v", fullline, r)
			returnNode = C4goErrorNode{}
		}
	}()
	line := fullline

	// This is a special case. I'm not sure if it's a bug in the clang AST
	// dumper. It should have children.
	if line == "array filler" {
		return parseArrayFiller(line), nil
	}

	parts := strings.SplitN(line, " ", 2)
	nodeName := parts[0]

	// skip node name
	if len(parts) > 1 {
		line = parts[1]
	}

	switch nodeName {
	case "AlignedAttr":
		return parseAlignedAttr(line), nil
	case "AllocSizeAttr":
		return parseAllocSizeAttr(line), nil
	case "AlwaysInlineAttr":
		return parseAlwaysInlineAttr(line), nil
	case "ArraySubscriptExpr":
		return parseArraySubscriptExpr(line), nil
	case "AsmLabelAttr":
		return parseAsmLabelAttr(line), nil
	case "AvailabilityAttr":
		return parseAvailabilityAttr(line), nil
	case "BinaryOperator":
		return parseBinaryOperator(line), nil
	case "BlockCommandComment":
		return parseBlockCommandComment(line), nil
	case "BreakStmt":
		return parseBreakStmt(line), nil
	case "BuiltinType":
		return parseBuiltinType(line), nil
	case "CallExpr":
		return parseCallExpr(line), nil
	case "CaseStmt":
		return parseCaseStmt(line), nil
	case "CharacterLiteral":
		return parseCharacterLiteral(line), nil
	case "CompoundLiteralExpr":
		return parseCompoundLiteralExpr(line), nil
	case "CompoundStmt":
		return parseCompoundStmt(line), nil
	case "ConditionalOperator":
		return parseConditionalOperator(line), nil
	case "ConstAttr":
		return parseConstAttr(line), nil
	case "ConstantArrayType":
		return parseConstantArrayType(line), nil
	case "ContinueStmt":
		return parseContinueStmt(line), nil
	case "CompoundAssignOperator":
		return parseCompoundAssignOperator(line), nil
	case "CStyleCastExpr":
		return parseCStyleCastExpr(line), nil
	case "CXXMemberCallExpr":
		return parseCXXMemberCallExpr(line), nil
	case "CXXRecord":
		return parseCXXRecord(line), nil
	case "CXXRecordDecl":
		return parseCXXRecordDecl(line), nil
	case "DecayedType":
		return parseDecayedType(line), nil
	case "DeclRefExpr":
		return parseDeclRefExpr(line), nil
	case "DeclStmt":
		return parseDeclStmt(line), nil
	case "DefaultStmt":
		return parseDefaultStmt(line), nil
	case "DeprecatedAttr":
		return parseDeprecatedAttr(line), nil
	case "DisableTailCallsAttr":
		return parseDisableTailCallsAttr(line), nil
	case "DoStmt":
		return parseDoStmt(line), nil
	case "ElaboratedType":
		return parseElaboratedType(line), nil
	case "EmptyDecl":
		return parseEmptyDecl(line), nil
	case "Enum":
		return parseEnum(line), nil
	case "EnumConstantDecl":
		return parseEnumConstantDecl(line), nil
	case "EnumDecl":
		return parseEnumDecl(line), nil
	case "EnumType":
		return parseEnumType(line), nil
	case "Field":
		return parseField(line), nil
	case "FieldDecl":
		return parseFieldDecl(line), nil
	case "FloatingLiteral":
		return parseFloatingLiteral(line), nil
	case "FormatArgAttr":
		return parseFormatArgAttr(line), nil
	case "FormatAttr":
		return parseFormatAttr(line), nil
	case "FunctionDecl":
		return parseFunctionDecl(line), nil
	case "FullComment":
		return parseFullComment(line), nil
	case "FunctionProtoType":
		return parseFunctionProtoType(line), nil
	case "ForStmt":
		return parseForStmt(line), nil
	case "HTMLStartTagComment":
		return parseHTMLStartTagComment(line), nil
	case "HTMLEndTagComment":
		return parseHTMLEndTagComment(line), nil
	case "GCCAsmStmt":
		return parseGCCAsmStmt(line), nil
	case "GotoStmt":
		return parseGotoStmt(line), nil
	case "IfStmt":
		return parseIfStmt(line), nil
	case "ImplicitCastExpr":
		return parseImplicitCastExpr(line), nil
	case "ImplicitValueInitExpr":
		return parseImplicitValueInitExpr(line), nil
	case "IncompleteArrayType":
		return parseIncompleteArrayType(line), nil
	case "IndirectFieldDecl":
		return parseIndirectFieldDecl(line), nil
	case "InitListExpr":
		return parseInitListExpr(line), nil
	case "InlineCommandComment":
		return parseInlineCommandComment(line), nil
	case "IntegerLiteral":
		return parseIntegerLiteral(line), nil
	case "LabelStmt":
		return parseLabelStmt(line), nil
	case "MallocAttr":
		return parseMallocAttr(line), nil
	case "MaxFieldAlignmentAttr":
		return parseMaxFieldAlignmentAttr(line), nil
	case "MemberExpr":
		return parseMemberExpr(line), nil
	case "ModeAttr":
		return parseModeAttr(line), nil
	case "NoInlineAttr":
		return parseNoInlineAttr(line), nil
	case "NoThrowAttr":
		return parseNoThrowAttr(line), nil
	case "NonNullAttr":
		return parseNonNullAttr(line), nil
	case "OffsetOfExpr":
		return parseOffsetOfExpr(line), nil
	case "PackedAttr":
		return parsePackedAttr(line), nil
	case "ParagraphComment":
		return parseParagraphComment(line), nil
	case "ParamCommandComment":
		return parseParamCommandComment(line), nil
	case "ParenExpr":
		return parseParenExpr(line), nil
	case "ParenType":
		return parseParenType(line), nil
	case "ParmVarDecl":
		return parseParmVarDecl(line), nil
	case "PointerType":
		return parsePointerType(line), nil
	case "PredefinedExpr":
		return parsePredefinedExpr(line), nil
	case "PureAttr":
		return parsePureAttr(line), nil
	case "QualType":
		return parseQualType(line), nil
	case "Record":
		return parseRecord(line), nil
	case "RecordDecl":
		return parseRecordDecl(line), nil
	case "RecordType":
		return parseRecordType(line), nil
	case "RestrictAttr":
		return parseRestrictAttr(line), nil
	case "ReturnStmt":
		return parseReturnStmt(line), nil
	case "ReturnsTwiceAttr":
		return parseReturnsTwiceAttr(line), nil
	case "SentinelAttr":
		return parseSentinelAttr(line), nil
	case "StmtExpr":
		return parseStmtExpr(line), nil
	case "StringLiteral":
		return parseStringLiteral(line), nil
	case "SwitchStmt":
		return parseSwitchStmt(line), nil
	case "TextComment":
		return parseTextComment(line), nil
	case "TranslationUnitDecl":
		return parseTranslationUnitDecl(line), nil
	case "TransparentUnionAttr":
		return parseTransparentUnionAttr(line), nil
	case "Typedef":
		return parseTypedef(line), nil
	case "TypedefDecl":
		return parseTypedefDecl(line), nil
	case "TypedefType":
		return parseTypedefType(line), nil
	case "UnaryExprOrTypeTraitExpr":
		return parseUnaryExprOrTypeTraitExpr(line), nil
	case "UnaryOperator":
		return parseUnaryOperator(line), nil
	case "UnusedAttr":
		return parseUnusedAttr(line), nil
	case "VAArgExpr":
		return parseVAArgExpr(line), nil
	case "VarDecl":
		return parseVarDecl(line), nil
	case "VerbatimBlockComment":
		return parseVerbatimBlockComment(line), nil
	case "VerbatimBlockLineComment":
		return parseVerbatimBlockLineComment(line), nil
	case "VerbatimLineComment":
		return parseVerbatimLineComment(line), nil
	case "VisibilityAttr":
		return parseVisibilityAttr(line), nil
	case "WarnUnusedResultAttr":
		return parseWarnUnusedResultAttr(line), nil
	case "WeakAttr":
		return parseWeakAttr(line), nil
	case "WhileStmt":
		return parseWhileStmt(line), nil
	case "NullStmt":
		return nil, nil
	}
	return C4goErrorNode{}, fmt.Errorf("unknown node type: `%v`", fullline)
}

func groupsFromRegex(rx, line string) map[string]string {
	// We remove tabs and newlines from the regex. This is purely cosmetic,
	// as the regex input can be quite long and it's nice for the caller to
	// be able to format it in a more readable way.
	fullRegexp := "^(?P<address>[0-9a-fx]+) " +
		strings.Replace(strings.Replace(rx, "\n", "", -1), "\t", "", -1)
	rx = fullRegexp + "[\\s]*$"

	re := util.GetRegex(rx)

	match := re.FindStringSubmatch(line)
	if len(match) == 0 {
		panic("could not match regexp with string\n" + rx + "\n" + line + "\n")
	}

	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}

	return result
}
