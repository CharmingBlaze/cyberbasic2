// Package bindings wires all foreign APIs onto a VM in one place (DRY, documented order).
package bindings

import (
	"cyberbasic/compiler/bindings/aisys"
	"cyberbasic/compiler/bindings/assets"
	"cyberbasic/compiler/bindings/audiosys"
	"cyberbasic/compiler/bindings/box2d"
	"cyberbasic/compiler/bindings/bullet"
	"cyberbasic/compiler/bindings/cameradot"
	"cyberbasic/compiler/bindings/dbp"
	"cyberbasic/compiler/bindings/ecs"
	"cyberbasic/compiler/bindings/effect"
	"cyberbasic/compiler/bindings/engine"
	"cyberbasic/compiler/bindings/game"
	"cyberbasic/compiler/bindings/indoor"
	"cyberbasic/compiler/bindings/inputmap"
	"cyberbasic/compiler/bindings/nakama"
	"cyberbasic/compiler/bindings/navigation"
	"cyberbasic/compiler/bindings/net"
	"cyberbasic/compiler/bindings/objects"
	"cyberbasic/compiler/bindings/physics2d"
	"cyberbasic/compiler/bindings/procedural"
	"cyberbasic/compiler/bindings/raylib"
	"cyberbasic/compiler/bindings/scene"
	"cyberbasic/compiler/bindings/shadersys"
	"cyberbasic/compiler/bindings/sql"
	"cyberbasic/compiler/bindings/std"
	"cyberbasic/compiler/bindings/tween"
	"cyberbasic/compiler/bindings/terrain"
	"cyberbasic/compiler/bindings/vegetation"
	"cyberbasic/compiler/bindings/water"
	"cyberbasic/compiler/bindings/windowdot"
	"cyberbasic/compiler/bindings/world"
	"cyberbasic/compiler/runtime"
	"cyberbasic/compiler/runtime/renderer"
	"cyberbasic/compiler/vm"
)

// RegisterOptions configures RegisterAll. Source is used for physics2d explicit-window detection.
type RegisterOptions struct {
	// Source is full program source (or accumulated REPL session). Used with runtime.DetectWindowMode.
	Source string
	// SkipRaylib skips raylib + flush override + renderer global hooks (for headless/unit tests).
	SkipRaylib bool
}

// RegisterAll installs every foreign binding on v in a fixed order.
//
// Order rationale:
//  1. Raylib core + flush override — foundation for drawing and EndFrame batching.
//  2. DBP runtime + renderer hooks — 3D/2D scene draw bridges into raylib.
//  3. Bullet, Box2D, high-level physics2d — physics before game/ecs layers that may depend on bodies.
//  4. Net, Nakama, Scene, Game — multiplayer and scene graph.
//  5. DBP 2D overlay, SQL, terrain stack, objects, procedural, water, vegetation, world, nav, indoor.
//  6. Std + v2 modules (audio, input, assets, shader, effect, camera.fx, tween, AI) + WINDOW + engine composition last.
//
// DBP terrain/water/object overlays must run after their native packages so integer-ID commands take precedence where intended.
func RegisterAll(v *vm.VM, opts RegisterOptions) error {
	if !opts.SkipRaylib {
		raylib.RegisterRaylib(v)
		runtime.RegisterFlushOverride(v)
	}
	dbp.RegisterDBP(v)
	if !opts.SkipRaylib {
		renderer.SetDraw3D(dbp.DrawScene3D)
		renderer.SetPreDraw2D(dbp.UpdateSpriteAnimations)
		renderer.SetVM(v)
	}
	bullet.RegisterBullet(v)
	box2d.RegisterBox2D(v)
	physics2d.RegisterPhysics2DHigh(v)
	ecs.RegisterECS(v)
	net.RegisterNet(v)
	nakama.RegisterNakama(v)
	scene.RegisterScene(v)
	game.RegisterGame(v)
	dbp.Register2D(v)
	sql.RegisterSQL(v)
	terrain.RegisterTerrain(v)
	dbp.RegisterTerrain(v)
	objects.RegisterObjects(v)
	dbp.RegisterDrawObjectOverlay(v)
	procedural.RegisterProcedural(v)
	water.RegisterWater(v)
	dbp.RegisterWater(v)
	vegetation.RegisterVegetation(v)
	world.RegisterWorld(v)
	navigation.RegisterNavigation(v)
	indoor.RegisterIndoor(v)
	std.RegisterStd(v)
	audiosys.RegisterAudiosys(v)
	inputmap.RegisterInputmap(v)
	assets.RegisterAssets(v)
	shadersys.RegisterShaderSys(v)
	effect.RegisterEffect(v)
	cameradot.RegisterCameraDot(v)
	tween.RegisterTween(v)
	aisys.RegisterAisys(v)
	windowdot.RegisterWindowDot(v)
	engine.RegisterEngine(v)

	physics2d.WorldEnsured = false
	physics2d.RequireExplicitWorld = runtime.DetectWindowMode(opts.Source) == runtime.ModeExplicit
	return nil
}
