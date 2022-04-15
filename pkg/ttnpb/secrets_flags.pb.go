// Code generated by protoc-gen-go-flags. DO NOT EDIT.
// versions:
// - protoc-gen-go-flags v0.0.0-dev
// - protoc              v3.17.3
// source: lorawan-stack/api/secrets.proto

package ttnpb

import (
	flagsplugin "github.com/TheThingsIndustries/protoc-gen-go-flags/flagsplugin"
	pflag "github.com/spf13/pflag"
)

// AddSelectFlagsForSecret adds flags to select fields in Secret.
func AddSelectFlagsForSecret(flags *pflag.FlagSet, prefix string, hidden bool) {
	flags.AddFlag(flagsplugin.NewBoolFlag(flagsplugin.Prefix("key-id", prefix), flagsplugin.SelectDesc(flagsplugin.Prefix("key-id", prefix), false), flagsplugin.WithHidden(hidden)))
	flags.AddFlag(flagsplugin.NewBoolFlag(flagsplugin.Prefix("value", prefix), flagsplugin.SelectDesc(flagsplugin.Prefix("value", prefix), false), flagsplugin.WithHidden(hidden)))
}

// SelectFromFlags outputs the fieldmask paths forSecret message from select flags.
func PathsFromSelectFlagsForSecret(flags *pflag.FlagSet, prefix string) (paths []string, err error) {
	if val, selected, err := flagsplugin.GetBool(flags, flagsplugin.Prefix("key_id", prefix)); err != nil {
		return nil, err
	} else if selected && val {
		paths = append(paths, flagsplugin.Prefix("key_id", prefix))
	}
	if val, selected, err := flagsplugin.GetBool(flags, flagsplugin.Prefix("value", prefix)); err != nil {
		return nil, err
	} else if selected && val {
		paths = append(paths, flagsplugin.Prefix("value", prefix))
	}
	return paths, nil
}

// AddSetFlagsForSecret adds flags to select fields in Secret.
func AddSetFlagsForSecret(flags *pflag.FlagSet, prefix string, hidden bool) {
	flags.AddFlag(flagsplugin.NewStringFlag(flagsplugin.Prefix("key-id", prefix), "", flagsplugin.WithHidden(hidden)))
	flags.AddFlag(flagsplugin.NewHexBytesFlag(flagsplugin.Prefix("value", prefix), "", flagsplugin.WithHidden(hidden)))
}

// SetFromFlags sets the Secret message from flags.
func (m *Secret) SetFromFlags(flags *pflag.FlagSet, prefix string) (paths []string, err error) {
	if val, changed, err := flagsplugin.GetString(flags, flagsplugin.Prefix("key_id", prefix)); err != nil {
		return nil, err
	} else if changed {
		m.KeyId = val
		paths = append(paths, flagsplugin.Prefix("key_id", prefix))
	}
	if val, changed, err := flagsplugin.GetBytes(flags, flagsplugin.Prefix("value", prefix)); err != nil {
		return nil, err
	} else if changed {
		m.Value = val
		paths = append(paths, flagsplugin.Prefix("value", prefix))
	}
	return paths, nil
}
