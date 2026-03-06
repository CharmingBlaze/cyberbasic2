package terrain

import (
	"fmt"
	"sync"

	"cyberbasic/compiler/vm"
)

// TerrainState holds mesh, material, and heightmap refs for a terrain instance.
type TerrainState struct {
	HeightmapID      string
	MeshID           string
	MaterialID       string
	SizeX            float32
	SizeZ            float32
	HeightScale      float32
	LODLevel         int
	CollisionEnabled bool
	Friction         float32
	Bounce           float32
	PosX             float32
	PosY             float32
	PosZ             float32
	Layers           [4]string // Texture paths/ids for splat layers 0-3
	SplatmapPath     string    // RGBA splatmap; if missing, use layer 0 only
	Visible          bool      // Visibility for Hide/Show
}

var (
	terrains   = make(map[string]*TerrainState)
	terrainSeq int
	terrainMu  sync.Mutex
	defaultMat string
)

// MakeTerrainFlat creates a flat terrain (plane mesh) without heightmap. Returns terrain id.
func MakeTerrainFlat(v *vm.VM, sizeX, sizeZ float32) (string, error) {
	if sizeX <= 0 {
		sizeX = 100
	}
	if sizeZ <= 0 {
		sizeZ = 100
	}
	resX, resZ := int32(32), int32(32)
	if sizeX > 100 {
		resX = int32(sizeX / 4)
	}
	if sizeZ > 100 {
		resZ = int32(sizeZ / 4)
	}
	if resX < 2 {
		resX = 2
	}
	if resZ < 2 {
		resZ = 2
	}
	meshRes, err := v.CallForeign("GenMeshPlane", []interface{}{sizeX, sizeZ, resX, resZ})
	if err != nil {
		return "", err
	}
	meshID, ok := meshRes.(string)
	if !ok || meshID == "" {
		return "", fmt.Errorf("GenMeshPlane did not return mesh id")
	}
	// Create flat heightmap for compatibility (all zeros)
	hmID, _ := createFlatHeightmap(int(sizeX), int(sizeZ))
	terrainMu.Lock()
	terrainSeq++
	id := fmt.Sprintf("terrain_%d", terrainSeq)
	terrains[id] = &TerrainState{
		HeightmapID: hmID,
		MeshID:      meshID,
		MaterialID:  "",
		SizeX:       sizeX,
		SizeZ:       sizeZ,
		HeightScale: 0,
		LODLevel:    0,
		Visible:     true,
	}
	terrainMu.Unlock()
	return id, nil
}

// TerrainCreate runs the mesh generator and registers a terrain. Returns terrain id.
func TerrainCreate(v *vm.VM, heightmapID string, sizeX, sizeZ, heightScale float32) (string, error) {
	meshID, err := GenTerrainMesh(v, heightmapID, sizeX, sizeZ, heightScale, 0)
	if err != nil {
		return "", err
	}
	terrainMu.Lock()
	terrainSeq++
	id := fmt.Sprintf("terrain_%d", terrainSeq)
	terrains[id] = &TerrainState{
		HeightmapID: heightmapID,
		MeshID:      meshID,
		MaterialID:  "",
		SizeX:       sizeX,
		SizeZ:       sizeZ,
		HeightScale: heightScale,
		LODLevel:    0,
		Visible:     true,
	}
	terrainMu.Unlock()
	return id, nil
}

// TerrainUpdate rebuilds the mesh from the current heightmap.
func TerrainUpdate(v *vm.VM, terrainID string) error {
	terrainMu.Lock()
	ts, ok := terrains[terrainID]
	terrainMu.Unlock()
	if !ok {
		return fmt.Errorf("unknown terrain id: %s", terrainID)
	}
	meshID, err := GenTerrainMesh(v, ts.HeightmapID, ts.SizeX, ts.SizeZ, ts.HeightScale, ts.LODLevel)
	if err != nil {
		return err
	}
	terrainMu.Lock()
	ts.MeshID = meshID
	terrainMu.Unlock()
	return nil
}

// DrawTerrain draws the terrain at the given position (calls DrawMesh).
func DrawTerrain(v *vm.VM, terrainID string, posX, posY, posZ float32) error {
	terrainMu.Lock()
	ts, ok := terrains[terrainID]
	terrainMu.Unlock()
	if !ok {
		return fmt.Errorf("unknown terrain id: %s", terrainID)
	}
	if !ts.Visible {
		return nil
	}
	matID := ts.MaterialID
	if matID == "" {
		if defaultMat == "" {
			res, err := v.CallForeign("LoadMaterialDefault", nil)
			if err != nil || res == nil {
				matID = "default"
			} else {
				defaultMat, _ = res.(string)
				matID = defaultMat
			}
		} else {
			matID = defaultMat
		}
	}
	_, err := v.CallForeign("DrawMesh", []interface{}{
		ts.MeshID, matID,
		posX, posY, posZ,
		float32(1), float32(1), float32(1),
	})
	return err
}

// SetTerrainMaterial sets the material id used when drawing.
func SetTerrainMaterial(terrainID, materialID string) error {
	terrainMu.Lock()
	defer terrainMu.Unlock()
	ts, ok := terrains[terrainID]
	if !ok {
		return fmt.Errorf("unknown terrain id: %s", terrainID)
	}
	ts.MaterialID = materialID
	return nil
}

// SetTerrainTexture is an alias that sets the terrain's material from a texture (user can create material first).
func SetTerrainTexture(terrainID, textureID string) error {
	return SetTerrainMaterial(terrainID, textureID)
}

// SetTerrainLOD sets LOD level; next TerrainUpdate will use it.
func SetTerrainLOD(terrainID string, lodLevel int) error {
	terrainMu.Lock()
	defer terrainMu.Unlock()
	ts, ok := terrains[terrainID]
	if !ok {
		return fmt.Errorf("unknown terrain id: %s", terrainID)
	}
	ts.LODLevel = lodLevel
	return nil
}

// GetTerrainState returns internal state for a terrain (for edit/query).
func GetTerrainState(terrainID string) *TerrainState {
	terrainMu.Lock()
	ts := terrains[terrainID]
	terrainMu.Unlock()
	return ts
}

// SetTerrainPosition stores position for DrawTerrain.
func SetTerrainPosition(terrainID string, posX, posY, posZ float32) error {
	terrainMu.Lock()
	defer terrainMu.Unlock()
	ts, ok := terrains[terrainID]
	if !ok {
		return fmt.Errorf("unknown terrain id: %s", terrainID)
	}
	ts.PosX, ts.PosY, ts.PosZ = posX, posY, posZ
	return nil
}

// SetTerrainLayer sets texture for layer 0-3 (for splatmap blending).
func SetTerrainLayer(terrainID string, layerIndex int, textureID string) error {
	terrainMu.Lock()
	defer terrainMu.Unlock()
	ts, ok := terrains[terrainID]
	if !ok {
		return fmt.Errorf("unknown terrain id: %s", terrainID)
	}
	if layerIndex >= 0 && layerIndex < 4 {
		ts.Layers[layerIndex] = textureID
	}
	return nil
}

// GenerateTerrainNoiseForTerrain replaces terrain heightmap with procedural noise. Deterministic via seed.
func GenerateTerrainNoiseForTerrain(v *vm.VM, terrainID string, seed int64, octaves int, scale float64) error {
	ts := GetTerrainState(terrainID)
	if ts == nil {
		return fmt.Errorf("unknown terrain id: %s", terrainID)
	}
	width := int(ts.SizeX)
	depth := int(ts.SizeZ)
	if width <= 0 {
		width = 64
	}
	if depth <= 0 {
		depth = 64
	}
	hmID, err := GenHeightmapNoise(width, depth, seed, octaves, scale)
	if err != nil {
		return err
	}
	terrainMu.Lock()
	ts.HeightmapID = hmID
	if ts.HeightScale <= 0 {
		ts.HeightScale = 20
	}
	terrainMu.Unlock()
	return TerrainUpdate(v, terrainID)
}

// SetTerrainSplatmap sets splatmap path; if missing, use layer 0 only.
func SetTerrainSplatmap(terrainID string, path string) error {
	terrainMu.Lock()
	defer terrainMu.Unlock()
	ts, ok := terrains[terrainID]
	if !ok {
		return fmt.Errorf("unknown terrain id: %s", terrainID)
	}
	ts.SplatmapPath = path
	return nil
}

// SetTerrainVisible sets visibility for Hide/Show.
func SetTerrainVisible(terrainID string, visible bool) {
	terrainMu.Lock()
	defer terrainMu.Unlock()
	if ts, ok := terrains[terrainID]; ok {
		ts.Visible = visible
	}
}

// TerrainDelete removes terrain from the registry and unloads mesh.
func TerrainDelete(v *vm.VM, terrainID string) error {
	terrainMu.Lock()
	ts, ok := terrains[terrainID]
	if !ok {
		terrainMu.Unlock()
		return nil
	}
	meshID := ts.MeshID
	delete(terrains, terrainID)
	terrainMu.Unlock()
	if meshID != "" {
		_, _ = v.CallForeign("UnloadMesh", []interface{}{meshID})
	}
	return nil
}

// TerrainClone creates a new terrain with the same state. Returns new terrain id.
func TerrainClone(v *vm.VM, srcID string) (string, error) {
	terrainMu.Lock()
	src, ok := terrains[srcID]
	if !ok {
		terrainMu.Unlock()
		return "", fmt.Errorf("unknown terrain id: %s", srcID)
	}
	terrainMu.Unlock()
	hmID, err := CloneHeightmap(src.HeightmapID)
	if err != nil {
		return "", err
	}
	meshID, err := GenTerrainMesh(v, hmID, src.SizeX, src.SizeZ, src.HeightScale, src.LODLevel)
	if err != nil {
		return "", err
	}
	terrainMu.Lock()
	terrainSeq++
	newID := fmt.Sprintf("terrain_%d", terrainSeq)
	terrains[newID] = &TerrainState{
		HeightmapID:      hmID,
		MeshID:           meshID,
		MaterialID:       src.MaterialID,
		SizeX:            src.SizeX,
		SizeZ:            src.SizeZ,
		HeightScale:      src.HeightScale,
		LODLevel:         src.LODLevel,
		CollisionEnabled: src.CollisionEnabled,
		Friction:          src.Friction,
		Bounce:            src.Bounce,
		PosX:             src.PosX,
		PosY:             src.PosY,
		PosZ:             src.PosZ,
		Layers:           src.Layers,
		SplatmapPath:     src.SplatmapPath,
		Visible:          src.Visible,
	}
	terrainMu.Unlock()
	return newID, nil
}
